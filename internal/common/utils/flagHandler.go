package utils

import (
	"flag"
	"os"
)

// HandleFlag обработчик флагов
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

	flag.Func("s", "Enable https flag", func(sFlagValue string) error {
		return os.Setenv("ENABLE_HTTPS", sFlagValue)
	})

	flag.Func("c", "Path of file config", func(cFlagValue string) error {
		return os.Setenv("CONFIG", cFlagValue)
	})
}
