package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/fileRepository"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/memoryRepository"
)

func CreateRepository() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	fmt.Println(filePath)
	if len(strings.TrimSpace(filePath)) == 0 {
		fmt.Println("in memory")
		return memoryRepository.NewURLMemoryRepository()
	} else {
		fmt.Println("in file")
		return fileRepository.NewURLFileRepository()
	}
}
