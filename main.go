package main

import (
	"fmt"
	"os"

	"github.com/dikeert/linkman/cmd"
	"github.com/dikeert/linkman/data"
)

const dataDir = "linkman"
const dataFile = "data.db"

func main() {
	if err := data.EnsureDataHome(dataDir); err == nil {
		if dataPath, err := data.GetFilePath(dataDir, dataFile); err == nil {
			cmd.Execute(dataPath, os.Args[1:])
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
	}

	os.Exit(3)
}
