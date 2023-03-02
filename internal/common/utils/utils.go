package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/indatabase"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/infile"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/inmemory"
)

// GetFilePath геттер env переменной FILE_STORAGE_PATH
func GetFilePath() string {
	return os.Getenv("FILE_STORAGE_PATH")
}

// GetServerAddress геттер env переменной SERVER_ADDRESS
func GetServerAddress() string {
	return os.Getenv("SERVER_ADDRESS")
}

// GetKey геттер env переменной SECRET_KEY
func GetKey() string {
	return os.Getenv("SECRET_KEY")
}

// GetDBAddress геттер env переменной DATABASE_DSN
func GetDBAddress() string {
	return os.Getenv("DATABASE_DSN")
}

// GetEnableHttps геттер env переменной ENABLE_HTTPS
func GetEnableHttps() bool {
	isEnableHttps, err := strconv.ParseBool(os.Getenv("ENABLE_HTTPS"))
	if err != nil {
		log.Fatal(err)
	}
	return isEnableHttps
}

// LoadEnvironments загрузка env переменных
func LoadEnvironments(envFilePath string) {
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// CreateMemoryOrFileRepository создание репозитория в зависимости от переменной
func CreateMemoryOrFileRepository(filePath string) repositories.URLRepository {
	if len(strings.TrimSpace(filePath)) != 0 {
		log.Println("in File")
		return infile.New(filePath)
	}

	log.Println("in Memory")
	return inmemory.New()
}

// CreateRepository создание репозитория в зависимости от переменной
func CreateRepository(filePath, dbAddress string) (repositories.URLRepository, error) {
	var repo repositories.URLRepository
	var err error
	if len(strings.TrimSpace(dbAddress)) != 0 {
		log.Println("in dataBase")
		repo, err = indatabase.New(dbAddress)
		if err != nil {
			return nil, err
		}
	} else {
		repo = CreateMemoryOrFileRepository(filePath)
	}
	return repo, nil
}
