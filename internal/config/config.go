package config

type Config struct {
	Host     string
	FilePath string
	Key      string
}

func New(serverAddress, filePath, key string) Config {
	return Config{
		Host:     serverAddress,
		FilePath: filePath,
		Key:      key,
	}
}
