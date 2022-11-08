package main

import (
	"fmt"
	"os"
)

func login() {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf(name)
}
