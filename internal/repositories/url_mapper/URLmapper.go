package url_mapper

import (
	"fmt"
	"os"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

func FromMapToSliceOfFullURL(userURLMap map[string]string) []entity.FullURL {
	fullURLsList := make([]entity.FullURL, 0, len(userURLMap))

	for short, original := range userURLMap {
		fullURL := entity.FullURL{
			ShortURL:    fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), short),
			OriginalURL: original,
		}
		fullURLsList = append(fullURLsList, fullURL)
	}

	return fullURLsList
}
