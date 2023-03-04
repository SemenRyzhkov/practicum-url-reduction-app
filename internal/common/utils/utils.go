package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
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

// GetEnableHTTPS геттер env переменной ENABLE_HTTPS
func GetEnableHTTPS() bool {
	isEnableHTTPS, err := strconv.ParseBool(os.Getenv("ENABLE_HTTPS"))
	if err != nil {
		log.Fatal(err)
	}
	return isEnableHTTPS
}

func GetConfigFilePath() string {
	return os.Getenv("CONFIG")
}

// LoadEnvironments загрузка env переменных
func LoadEnvironments(envFilePath string) {
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// CreateConfig создание конфига
func CreateConfig(
	serverAddress,
	filePath,
	key,
	dbAddress,
	configFilePath string,
	enableHTTPS bool,
) (config.Config, error) {
	if environmentsIsEmpty(serverAddress, filePath, key, dbAddress, enableHTTPS) {
		return config.LoadConfiguration(configFilePath)
	} else {
		return config.New(serverAddress, filePath, key, dbAddress, enableHTTPS), nil
	}
}

func environmentsIsEmpty(serverAddress, filePath, key, dbAddress string, enableHTTPS bool) bool {
	return len(serverAddress) == 0 && len(filePath) == 0 && len(key) == 0 && len(dbAddress) == 0 && !enableHTTPS
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
