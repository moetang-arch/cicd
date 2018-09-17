package main

import "github.com/moetang-arch/cicd/setup/setupcontroller/goproxy"

func main() {
	goproxy.StartSync(30132)

	// planning
	// get source
	// 1. `GOPATH=SHARED_PATH go mod tidy` with go.mod of the source
	// 2. serve HTTP_SERVICE in SHARED_PATH as GOPROXY
	// 3. run `GOCACHE=off GOPATH=TEMP_PATH GOPROXY=HTTP_SERVICE go mod tidy` in the source folder
	// 4. external: invoke testing,building,packaging
	// 5. clean up environment: delete TEMP_PATH
}
