package config

import (
	"encoding/json"
	"os"
)

// Config конфигкрация приложения
type Config struct {
	Host            string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FilePath        string `json:"file_storage_path"`
	Key             string `json:"secret_key"`
	DataBaseAddress string `json:"database_dsn"`
	TrustedSubnet   string `json:"trusted_subnet"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// New конструктор Config
func New(serverAddress, filePath, key, dbAddress string, enableHTTPS bool) Config {
	return Config{
		Host:            serverAddress,
		FilePath:        filePath,
		Key:             key,
		DataBaseAddress: dbAddress,
		EnableHTTPS:     enableHTTPS,
	}
}

// LoadConfiguration загрузка Config из файла
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
