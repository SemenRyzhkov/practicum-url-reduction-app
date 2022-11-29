package utils

import (
	"fmt"
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

func GetKey() string {
	return os.Getenv("SECRET_KEY")
}

func GetDBAddress() string {
	return os.Getenv("DATABASE_DSN")
}

func LoadEnvironments(envFilePath string) {
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func CreateRepository(filePath string) repositories.URLRepository {
	if len(strings.TrimSpace(filePath)) != 0 {
		fmt.Println("in File")

		return infile.New(filePath)
	}
	fmt.Println("in Memory")
	return inmemory.New()
}
