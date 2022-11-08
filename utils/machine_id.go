package utils

import (
	"log"

	"github.com/denisbrodbeck/machineid"
)

func getMachineID() string {
	id, err := machineid.ProtectedID("file-sync")
	if err != nil {
		log.Fatal(err)
	}
	return getSha1Str(id)
}
