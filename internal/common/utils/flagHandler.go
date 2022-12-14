package utils

import (
	"flag"
	"os"
)

func HandleFlag() {
	flag.Func("a", "HTTP server address", func(aFlagValue string) error {
		return os.Setenv("SERVER_ADDRESS", aFlagValue)
	})

	flag.Func("b", "Base Url", func(bFlagValue string) error {
		return os.Setenv("BASE_URL", bFlagValue)
	})

	flag.Func("f", "Path of file storage", func(fFlagValue string) error {
		return os.Setenv("FILE_STORAGE_PATH", fFlagValue)
	})

	flag.Func("d", "Address of db connection", func(dFlagValue string) error {
		return os.Setenv("DATABASE_DSN", dFlagValue)
	})
}
