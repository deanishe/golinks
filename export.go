package main

import (
	"fmt"

	"github.com/pelletier/go-toml"
	"github.com/prologic/bitcask"
)

// Export legacy Bitcask database to TOML.
func exportDatabase(dbpath string) error {
	var (
		bookmarks = map[string]string{}
		db        *bitcask.Bitcask
		data      []byte
		err       error
	)

	if db, err = bitcask.Open(dbpath); err != nil {
		return err
	}
	defer db.Close()

	err = db.Scan([]byte("bookmark_"), func(key []byte) error {
		val, err := db.Get(key)
		if err != nil {
			return err
		}

		bookmarks[string(key)[9:]] = string(val)
		return nil
	})

	if err != nil {
		return err
	}

	if data, err = toml.Marshal(bookmarks); err != nil {
		return err
	}

	fmt.Print(string(data))

	return nil
}
