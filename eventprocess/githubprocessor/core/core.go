package core

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
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
	secretKey = "123456"
	debug     = true

	ps = NewPushService()

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/events/github", func(writer http.ResponseWriter, request *http.Request) {
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		sig := request.Header.Get(GITHUB_HEADER_SIGNATURE)
		mac := hmac.New(sha1.New, []byte(secretKey))
		mac.Write(data)
		expectedMAC := mac.Sum(nil)
		actualMAC := fmt.Sprintf("sha1=%x", expectedMAC)
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
	})

	http.ListenAndServe("0.0.0.0:4567", mux)
}
