package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

var (
	db *badger.DB
)

func dbInit() *badger.DB {
	loggerLevel := badger.ERROR

	if isDev {
		loggerLevel = badger.INFO
	}
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions(getDBPath()).WithLoggingLevel(loggerLevel)
	opts.IndexCacheSize = 100 << 20 // 100 mb or some other size based on the amount of data
	// opts.EncryptionKey = []byte("badger")

	db, err := badger.Open(opts)

	if err != nil {
		log.Fatal("DB Error: " + err.Error())
	}

	return db
}
