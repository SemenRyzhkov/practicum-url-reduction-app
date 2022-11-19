package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/infile"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/inmemory"
)

func CreateRepository() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	fmt.Println(filePath)
	if len(strings.TrimSpace(filePath)) == 0 {
		fmt.Println("in memory")
		return inmemory.New()
	} else {
		fmt.Println("in file")
		return infile.New()
	}
}
