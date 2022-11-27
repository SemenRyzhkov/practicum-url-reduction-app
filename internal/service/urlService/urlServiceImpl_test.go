package urlService

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/testutils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

func Test_urlServiceImpl_GetUrlById(t *testing.T) {
	utils.LoadEnvironments("../../../.env")
	tests := []struct {
		repo    repositories.URLRepository
		name    string
		urlID   string
		want    string
		wantErr bool
		userID  string
	}{
		{
			repo:    utils.CreateRepository(utils.GetFilePath()),
			name:    "positive test #1",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d28",
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath()),
			name:    "not found test #2",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d29",
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLService(tt.repo)
			u.ReduceAndSaveURL(tt.userID, tt.want)

			got, err := u.GetURLByID(tt.userID, tt.urlID)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
		testutils.AfterTest()
	}
}

func Test_urlServiceImpl_ReduceAndSaveUrl(t *testing.T) {
	utils.LoadEnvironments("../../../.env")
	tests := []struct {
		repo    repositories.URLRepository
		name    string
		saveURL string
		want    string
		wantErr bool
		userID  string
	}{
		{
			repo:    utils.CreateRepository(utils.GetFilePath()),
			name:    "positive test #1",
			saveURL: "yandex1.com",
			want:    "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b",
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath()),
			name:    "duplicate test #2",
			saveURL: "yandex.com",
			want:    "http://localhost:8080/XVlBz",
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLService(tt.repo)
			got, _ := u.ReduceAndSaveURL(tt.userID, tt.saveURL)
			if tt.wantErr {
				_, err := u.ReduceAndSaveURL(tt.userID, tt.saveURL)
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
		testutils.AfterTest()
	}
}

func Test_urlServiceImpl_ReduceUrlToJSON(t *testing.T) {
	utils.LoadEnvironments("../../../.env")
	tests := []struct {
		repo    repositories.URLRepository
		name    string
		request entity.URLRequest
		want    entity.URLResponse
		wantErr bool
		userID  string
	}{
		{
			repo:    utils.CreateRepository(utils.GetFilePath()),
			name:    "reducing JSON test #1",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex1.com"},
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath()),
			name:    "duplicate test #2",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex1.com"},
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewURLService(tt.repo)
			got, _ := u.ReduceURLToJSON(tt.userID, tt.request)
			if tt.wantErr {
				_, err := u.ReduceURLToJSON(tt.userID, tt.request)
				assert.NotNil(t, err)
			} else {
				assert.Equalf(t, tt.want, got, "ReduceURLToJSON(%v)", tt.request)
			}
		})
		testutils.AfterTest()
	}
}
