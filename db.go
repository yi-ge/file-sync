package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

var (
	db *badger.DB
)

func dbInit() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(getDBPath())
	opts.IndexCacheSize = 100 << 20 // 100 mb or some other size based on the amount of data
	// opts.EncryptionKey = []byte("badger")

	db, err := badger.Open(opts)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Your code here…
}
