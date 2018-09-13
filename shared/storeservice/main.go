package main

import (
	"flag"

	"github.com/moetang-arch/cicd/shared/storeservice/internal/core"
)

var (
	STORE_PATH string
)

func init() {
	flag.StringVar(&STORE_PATH, "path", "store.db", "store file path")
}

func main() {
	flag.Parse()

	storeService, err := core.NewStoreService(STORE_PATH)
	if err != nil {
		panic(err)
	}
	err = storeService.Init(func(helperFactory core.HelperFactory) (err error) {
		defer func() {
			r := recover()
			if e, ok := r.(error); ok {
				err = e
			}
		}()
		helper := helperFactory.GetHelper()
		createTableFunc(helper.CreateTable, TableTemp)
		createTableFunc(helper.CreateTable, TableRepository)
	})
	if err != nil {
		panic(err)
	}

	defer storeService.Close()
}

func createTableFunc(fn func(name string) error, name string) {
	err := fn(name)
	if err != nil {
		panic(err)
	}
}
