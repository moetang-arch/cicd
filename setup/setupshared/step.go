package setupshared

import (
	"path/filepath"
	"os"
)

func PrepareRootFolder(root string) (string, error) {
	rootFolder := filepath.Join(root, RandomString(32))
	err := os.MkdirAll(rootFolder, 0755)
	if err != nil {
		return "", err
	}
	return rootFolder, nil
}
