package config

type Config struct {
	Host            string
	FilePath        string
	Key             string
	DataBaseAddress string
}

func New(serverAddress, filePath, key, dbAddress string) Config {
	return Config{
		Host:            serverAddress,
		FilePath:        filePath,
		Key:             key,
		DataBaseAddress: dbAddress,
	}
}
