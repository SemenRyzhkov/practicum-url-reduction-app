package testutils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnvironments() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env filerepository")
	}
}

func AfterTest() {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	fmt.Println(filePath)
	if len(strings.TrimSpace(filePath)) == 0 {
		return
	} else {
		e := os.Truncate(filePath, 0)
		if e != nil {
			log.Fatal(e)
		}
	}
}
