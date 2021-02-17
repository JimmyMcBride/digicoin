package utils

import "log"

// HandleErr takes in an `err` and handles it for you.
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
