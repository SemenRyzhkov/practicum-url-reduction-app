package config

type Config struct {
	Host     string
	FilePath string
}

func New(serverAddress, filePath string) Config {
	return Config{
		Host:     serverAddress,
		FilePath: filePath,
	}
}
