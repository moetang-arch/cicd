package goproxy

import (
	"net/http"
)

/*
	1. get level with
	a list of all known versions of the given module, one per line
	/????/????/@v/list
	/????/????/@latest
	/????/@v/list
	/????/@latest
	if 404 then clients request parent directory

	2. get info
	/????/@v/<version>.info
	content json:
	type Info struct {
    	Version string    // version string
    	Time    time.Time // commit time
	}

	3. go.mod
	/????/@v/<version>.mod
	returns the go.mod file for that version of the given module

	4. zip
	/????/@v/<version>.zip
	returns the zip archive for that version of the given module
	every file path in the archive must begin with <module>@<version>/

	5. replacing every uppercase letter with an exclamation mark
	   followed by the corresponding lower-case letter:
	   github.com/Azure encodes as github.com/!azure

	6. serving $GOPATH/pkg/mod/cache/download at (or copying it to)
	   https://example.com/proxy would let other users
	   access those cached module versions with GOPROXY=https://example.com/proxy.
*/
func StartSync(addr string, path string) {
	err := http.ListenAndServe(addr, http.FileServer(http.Dir(path)))
	if err != nil {
		panic(err)
	}
}
