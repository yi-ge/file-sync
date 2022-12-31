package utils

import (
	"log"

	"github.com/denisbrodbeck/machineid"
)

func GetMachineID() string {
	id, err := machineid.ProtectedID("file-sync")
	if err != nil {
		log.Fatal(err)
	}
	return GetSha1Str(id)
}

func GetMachineIDUseSHA256() string {
	id, err := machineid.ProtectedID("file-sync")
	if err != nil {
		log.Fatal(err)
	}
	return GetSha1Str(id)
}
