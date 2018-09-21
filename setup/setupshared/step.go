package setupshared

import (
	"path/filepath"
	"os"

	"github.com/moetang-arch/cicd/shared/model"
)

func PrepareRootFolder(root string) (string, error) {
	rootFolder := filepath.Join(root, RandomString(32))
	err := os.MkdirAll(rootFolder, 0755)
	if err != nil {
		return "", err
	}
	return rootFolder, nil
}

func CheckoutCode(root string, pe *model.PushEvent) (codepath string, err error) {
	//TODO
}

func LoadGoModFile(path string) (content string, err error) {
	//TODO
}
