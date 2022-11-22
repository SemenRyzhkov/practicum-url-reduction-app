package utils

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/infile"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/inmemory"
)

func GetFilePath() string {
	return os.Getenv("FILE_STORAGE_PATH")
}

func GetServerAddress() string {
	return os.Getenv("SERVER_ADDRESS")
}

func LoadEnvironments(envFilePath string) {
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func CreateRepository(filePath string) repositories.URLRepository {
	if len(strings.TrimSpace(filePath)) == 0 {
		return inmemory.New()
	} else {
		return infile.New(filePath)
	}
}
