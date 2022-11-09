package utils

import "log"

// check does error checking
func check(e error) {
	if e != nil {
		log.Printf("error: %v", e)
	}
}
