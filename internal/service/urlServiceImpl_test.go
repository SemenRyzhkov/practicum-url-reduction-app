package service

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/fileStorage"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/memoryStorage"
)

func testSetup() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env fileStorage")
	}
}

func afterTest() {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	fmt.Println(filePath)
	if len(strings.TrimSpace(filePath)) == 0 {
		return
	} else {
		e := os.Truncate(filePath, 0)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func createRepository() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	fmt.Println(filePath)
	if len(strings.TrimSpace(filePath)) == 0 {
		return memoryStorage.NewURLMemoryRepository()
	} else {
		return fileStorage.NewURLFileRepository()
	}
}

func Test_urlServiceImpl_GetUrlById(t *testing.T) {
	testSetup()
	tests := []struct {
		repo    repositories.URLRepository
		name    string
		urlID   string
		want    string
		wantErr bool
	}{
		{
			repo:    createRepository(),
			name:    "positive test #1",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d28",
			wantErr: false,
		},
		{
			repo:    createRepository(),
			name:    "not found test #2",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d29",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLService(tt.repo)
			u.ReduceAndSaveURL(tt.want)

			got, err := u.GetURLByID(tt.urlID)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
		afterTest()
	}
}

func Test_urlServiceImpl_ReduceAndSaveUrl(t *testing.T) {
	testSetup()
	tests := []struct {
		repo    repositories.URLRepository
		name    string
		saveURL string
		want    string
		wantErr bool
	}{
		{
			repo:    createRepository(),
			name:    "positive test #1",
			saveURL: "yandex1.com",
			want:    "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b",
			wantErr: false,
		},
		{
			repo:    createRepository(),
			name:    "duplicate test #2",
			saveURL: "yandex.com",
			want:    "http://localhost:8080/XVlBz",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLService(tt.repo)
			got, _ := u.ReduceAndSaveURL(tt.saveURL)
			if tt.wantErr {
				_, err := u.ReduceAndSaveURL(tt.saveURL)
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
		afterTest()

	}

}

func Test_urlServiceImpl_ReduceUrlToJSON(t *testing.T) {
	testSetup()
	tests := []struct {
		repo    repositories.URLRepository
		name    string
		request entity.URLRequest
		want    entity.URLResponse
		wantErr bool
	}{
		{
			repo:    createRepository(),
			name:    "reducing JSON test #1",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex1.com"},
			wantErr: false,
		},
		{
			repo:    createRepository(),
			name:    "duplicate test #2",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex1.com"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLService(tt.repo)
			got, _ := u.ReduceURLToJSON(tt.request)
			if tt.wantErr {
				_, err := u.ReduceURLToJSON(tt.request)
				assert.NotNil(t, err)
			} else {
				assert.Equalf(t, tt.want, got, "ReduceURLToJSON(%v)", tt.request)
			}
		})
		afterTest()
	}
}
