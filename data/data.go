package data

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
)

func EnsureDataHome(dir string) error {
	path := xdg.DataHome + dir
	if err := os.MkdirAll(path, 0700); err == nil {
		return nil
	} else {
		return fmt.Errorf("Unable to create data directory: %s", err)
	}
}

func GetFilePath(dir string, file string) (string, error) {
	if path, err := xdg.DataFile(dir + "/" + file); err == nil {
		return path, nil
	} else {
		return "", fmt.Errorf("Unable for find path to file [%s/%s]: %s",
			dir, file, err)
	}
}
