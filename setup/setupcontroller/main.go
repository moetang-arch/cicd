package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/moetang-arch/cicd/setup/setupcontroller/goproxy"
	"github.com/moetang-arch/cicd/setup/setupshared"
	"github.com/moetang-arch/cicd/shared/model"
)

var (
	// required
	TBP_AGENT_PATH string
	NSQ_ADDR       string
	SHARED_PATH    string
	HTTP_SERVICE   string
	GOPROXY_CONFIG string
	TEMP_PATH      string

	// optional
	GO_BIN_PATH string
)

var (
	logger *log.Logger
	queue  = make(chan *model.PushEvent)
)

func init() {
	flag.StringVar(&TBP_AGENT_PATH, "tbpagent", "", "testing/building/packaging agent binary path")
	flag.StringVar(&GO_BIN_PATH, "gobin", "", "go binary path")
	flag.StringVar(&NSQ_ADDR, "nsqaddr", "127.0.0.1:4150", "nsq address")
	flag.StringVar(&SHARED_PATH, "sharedpath", "", "shared path for GOPROXY storage")
	flag.StringVar(&HTTP_SERVICE, "httpservice", "0.0.0.0:31111", "http service for GOPROXY")
	flag.StringVar(&GOPROXY_CONFIG, "goproxy", "127.0.0.1:31111", "GOPROXY config")
	flag.StringVar(&TEMP_PATH, "temppath", "", "TEMP_PATH for requests")

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
	if GOPROXY_CONFIG == "" {
		fmt.Println("GOPROXY_CONFIG is empty.")
		os.Exit(1)
	}

	// task queue
	go envSetup()

	closer, err := startNsq(NSQ_ADDR, processor)
	if err != nil {
		fmt.Println("start nsq processor error.", err)
		os.Exit(1)
	}
	defer closer.Close()

	goproxy.StartSync(HTTP_SERVICE, filepath.Join(SHARED_PATH, "pkg", "mod", "cache", "download"))
}

func processor(data []byte) error {
	pe := new(model.PushEvent)
	err := pe.FromBytes(data)
	if err != nil {
		log.Println("unmarshal error.", err)
		return err
	}
	queue <- pe
	return nil
}

func envSetup() {
	for {
		pe := <-queue
		// 1. prepare directory
		rootFolder, err := setupshared.PrepareRootFolder(TEMP_PATH)
		if err != nil {
			logger.Println("prepare root folder error.", err)
			return
		}
		// 2. checkout code
		codePath, err := setupshared.CheckoutCode(rootFolder, pe)
		if err != nil {
			logger.Println("checkout code error.", err)
			return
		}
		// 3. try to load go.mod
		gomodContent, err := setupshared.LoadGoModFile(filepath.Join(codePath, "go.mod"))
		if err != nil {
			logger.Println("load go.mod file error.", err)
			return
		}
		if len(gomodContent) > 0 {
			//TODO 3.1. run `GOPATH=SHARED_PATH go mod download <module>` for each `require` element of `go mod edit -json go.mod`
			//TODO 3.2. run `GOCACHE=off GOPATH=TEMP_PATH GOPROXY=GOPROXY_CONFIG go mod tidy` in the source folder
		}
		//TODO 4. external: invoke testing,building,packaging
		//TODO 5. clean up environment: delete TEMP_PATH
	}
}
func checkoutCode(s string) {

}

func rollbackEvent() {
	//FIXME finish this function to support rollback event when error occurs
}
