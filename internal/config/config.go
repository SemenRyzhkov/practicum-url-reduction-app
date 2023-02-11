package config

// Config конфигкрация приложения
type Config struct {
	Host            string
	FilePath        string
	Key             string
	DataBaseAddress string
}

// New конструктор Config
func New(serverAddress, filePath, key, dbAddress string) Config {
	return Config{
		Host:            serverAddress,
		FilePath:        filePath,
		Key:             key,
		DataBaseAddress: dbAddress,
	}
}
