package utils

import (
	"os"
	"strings"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/infile"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/inmemory"
)

func CreateRepository() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if len(strings.TrimSpace(filePath)) == 0 {
		return inmemory.New()
	} else {
		return infile.New()
	}
}
