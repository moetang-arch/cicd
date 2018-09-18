package main

import (
	"flag"
	"log"
	"os"
	"fmt"

	"github.com/moetang-arch/cicd/setup/setupcontroller/goproxy"
	"github.com/moetang-arch/cicd/shared/model"
	"path/filepath"
)

var (
	// required
	TBP_AGENT_PATH string
	NSQ_ADDR       string
	SHARED_PATH    string
	HTTP_SERVICE   string

	// optional
	GO_BIN_PATH string
)

var (
	logger *log.Logger
)

func init() {
	flag.StringVar(&TBP_AGENT_PATH, "tbpagent", "", "testing/building/packaging agent binary path")
	flag.StringVar(&GO_BIN_PATH, "gobin", "", "go binary path")
	flag.StringVar(&NSQ_ADDR, "nsqaddr", "127.0.0.1:4150", "nsq address")
	flag.StringVar(&SHARED_PATH, "sharedpath", "", "shared path for GOPROXY storage")
	flag.StringVar(&HTTP_SERVICE, "httpservice", "0.0.0.0:31111", "http service for GOPROXY")

	logger = log.New(os.Stderr, "", log.LstdFlags)

	flag.Parse()
}

func main() {
	if TBP_AGENT_PATH == "" {
		fmt.Println("tbpagent is empty.")
		os.Exit(1)
	}
	if NSQ_ADDR == "" {
		fmt.Println("nsqaddr is empty.")
		os.Exit(1)
	}
	if SHARED_PATH == "" {
		fmt.Println("sharedpath is empty.")
		os.Exit(1)
	}
	if HTTP_SERVICE == "" {
		fmt.Println("htppservice is empty.")
		os.Exit(1)
	}

	closer, err := startNsq(NSQ_ADDR, processor)
	if err != nil {
		fmt.Println("start nsq processor error.", err)
		os.Exit(1)
	}
	defer closer.Close()

	goproxy.StartSync(HTTP_SERVICE, filepath.Join(SHARED_PATH, "pkg", "mod", "cache", "download"))

	// planning
	// get source & analyze module name then mv the source folder to it
	// 1. `GOPATH=SHARED_PATH go mod tidy` with go.mod of the source
	// 2. serve HTTP_SERVICE in SHARED_PATH as GOPROXY
	// 3. run `GOCACHE=off GOPATH=TEMP_PATH GOPROXY=HTTP_SERVICE go mod tidy` in the source folder
	// 4. external: invoke testing,building,packaging
	// 5. clean up environment: delete TEMP_PATH
}

func processor(data []byte) error {
	pe := new(model.PushEvent)
	err := pe.FromBytes(data)
	if err != nil {
		log.Println("unmarshal error.", err)
		return err
	}
	//TODO
	return nil
}
