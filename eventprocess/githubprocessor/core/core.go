package core

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	GITHUB_HEADER_EVENT_TYPE  = "X-GitHub-Event"
	GITHUB_HEADER_DELIVERY_ID = "X-GitHub-Delivery"
	GITHUB_HEADER_SIGNATURE   = "X-Hub-Signature"

	BRANCH_PREFIX = "refs/heads/"
)

var (
	secretKey  string
	debug      bool
	httpListen string

	nsqAddr string
	ps      *PushService

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func init() {
	flag.StringVar(&secretKey, "secretkey", "123456", "secret key for validating events")
	flag.BoolVar(&debug, "v", false, "enable verbose mode")
	flag.StringVar(&nsqAddr, "nsqaddr", "127.0.0.1:4150", "nsq address (required)")
	flag.StringVar(&httpListen, "listen", "0.0.0.0:4567", "http listen address")
}

func Start() {
	flag.Parse()

	ps = NewPushService(nsqAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/events/github", HttpHandler)

	http.ListenAndServe(httpListen, mux)
}

func HttpHandler(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sig := request.Header.Get(GITHUB_HEADER_SIGNATURE)
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write(data)
	expectedMAC := mac.Sum(nil)
	actualMAC := fmt.Sprintf("sha1=%x", expectedMAC) // github format
	checkResult := sig == actualMAC
	if !checkResult {
		logger.Println("signature not match. expected:", sig, "actual:", actualMAC)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	deliveryId := request.Header.Get(GITHUB_HEADER_DELIVERY_ID)
	eventType := request.Header.Get(GITHUB_HEADER_EVENT_TYPE)

	switch eventType {
	case "push":
		pushEvent := new(PushEvent)
		if err := json.Unmarshal(data, pushEvent); err != nil {
			log.Println("unmarshal data error. error:", err)
			if debug {
				log.Println("data:", string(data))
			}
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		pushEvent.DeliveryId = deliveryId

		if !strings.HasPrefix(pushEvent.Ref, BRANCH_PREFIX) {
			if debug {
				log.Println("ref is not branch. ref:", pushEvent.Ref)
			}
			writer.WriteHeader(http.StatusOK)
			return
		}

		pushEvent.Ref = pushEvent.Ref[len(BRANCH_PREFIX):]

		if err := ps.SendPushEvent(pushEvent); err != nil {
			log.Println("send push event to event bus error. error:", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if debug {
			log.Println("send push event of branch:", pushEvent.Ref, "commit:", pushEvent.AfterCommit, "clone_url:", pushEvent.Repository.CloneUrl)
		}
		writer.WriteHeader(http.StatusOK)
		return
	default:
		log.Println("skip unknown eventType:", eventType)
		if debug {
			log.Println("data:", string(data))
		}
		writer.WriteHeader(http.StatusOK)
		return
	}
}
