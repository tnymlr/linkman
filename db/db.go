package db

import (
	"fmt"

	"github.com/asdine/storm"
)

//Open tries open database located at specified path.
func Open(path string) (*storm.DB, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open database: %s", err)
	}

	return db, nil
}
