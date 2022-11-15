package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/filerepository"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/memoryrepository"
)

func CreateRepository() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	fmt.Println(filePath)
	if len(strings.TrimSpace(filePath)) == 0 {
		fmt.Println("in memory")
		return memoryrepository.NewURLMemoryRepository()
	} else {
		fmt.Println("in file")
		return filerepository.NewURLFileRepository()
	}
}
