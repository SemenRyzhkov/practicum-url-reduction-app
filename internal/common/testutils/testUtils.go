package testutils

import (
	"log"
	"os"
	"strings"
)

func AfterTest() {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if len(strings.TrimSpace(filePath)) == 0 {
		return
	} else {
		e := os.Truncate(filePath, 0)
		if e != nil {
			log.Fatal(e)
		}
	}
}
