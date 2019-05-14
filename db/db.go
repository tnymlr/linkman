package db

import (
	"fmt"

	"github.com/asdine/storm"
)

func Open(path string) (*storm.DB, error) {
	if db, err := storm.Open(path); err == nil {
		return db, nil
	} else {
		return nil, fmt.Errorf("Unable to open database: %s", err)
	}
}
