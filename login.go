package main

import (
	"fmt"
	"log"
	"os"

	"github.com/denisbrodbeck/machineid"
)

func getMachineID() string {
	id, err := machineid.ProtectedID("file-sync")
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func login() {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf(name)
}
