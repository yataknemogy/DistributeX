package common

import (
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func CreateDirectory(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		CheckError(err)
	}
}
