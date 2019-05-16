package data

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
)

//EnsureDataHome make sture that directory dir exists at
//location XDG_DATA_HOME.
func EnsureDataHome(dir string) error {
	path := xdg.DataHome + dir

	err := os.MkdirAll(path, 0700)
	if err != nil {
		return fmt.Errorf("Unable to create data directory: %s", err)
	}

	return nil
}

//GetFilePath returns path to the file in directory dir
//that is located in XDG_DATA_HOME.
func GetFilePath(dir string, file string) (string, error) {
	path, err := xdg.DataFile(dir + "/" + file)
	if err != nil {
		return "", fmt.Errorf(
			"Unable for find path to file [%s/%s]: %s",
			dir, file, err,
		)
	}

	return path, nil
}
