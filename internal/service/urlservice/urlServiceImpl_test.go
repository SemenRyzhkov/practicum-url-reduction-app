package urlservice

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/testutils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/mock_repositories"
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
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "positive test #1",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d28",
			wantErr: false,
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "not found test #2",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d29",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(strings.TrimSpace(utils.GetDBAddress())) != 0 {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				s := mock_repositories.NewMockURLRepository(ctrl)

				ctx := context.Background()

				if tt.wantErr {
					s.EXPECT().FindByID(ctx, tt.urlID).Return("", fmt.Errorf("error"))
					u := New(s)
					_, err := u.GetURLByID(ctx, tt.urlID)
					assert.NotNil(t, err)
				} else {
					s.EXPECT().FindByID(ctx, tt.urlID).Return("yandex.com", nil)
					u := New(s)
					got, _ := u.GetURLByID(ctx, tt.urlID)
					assert.Equal(t, tt.want, got)
				}
			} else {
				u := New(tt.repo)
				u.ReduceAndSaveURL(nil, tt.userID, tt.want)

				got, err := u.GetURLByID(nil, tt.urlID)
				if tt.wantErr {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, tt.want, got)
				}
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
		urlID   string
	}{
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "positive test #1",
			saveURL: "yandex1.com",
			want:    "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b",
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   "dc605989f530a3dfe9f7edacf1b3965b",
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "duplicate test #2",
			saveURL: "yandex.com",
			want:    "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b",
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   "31aa70fc8589c52a763a2df36f304d28",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(strings.TrimSpace(utils.GetDBAddress())) != 0 {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				s := mock_repositories.NewMockURLRepository(ctrl)

				ctx := context.Background()

				if tt.wantErr {
					s.EXPECT().Save(ctx, tt.userID, tt.urlID, tt.saveURL).Return(fmt.Errorf("error"))
					u := New(s)
					_, err := u.ReduceAndSaveURL(ctx, tt.userID, tt.saveURL)
					assert.NotNil(t, err)
				} else {
					s.EXPECT().Save(ctx, tt.userID, tt.urlID, tt.saveURL).Return(nil)
					u := New(s)
					got, _ := u.ReduceAndSaveURL(ctx, tt.userID, tt.saveURL)
					assert.Equal(t, tt.want, got)
				}
			} else {
				u := New(tt.repo)
				got, _ := u.ReduceAndSaveURL(nil, tt.userID, tt.saveURL)
				if tt.wantErr {
					_, err := u.ReduceAndSaveURL(nil, tt.userID, tt.saveURL)
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, tt.want, got)
				}
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
		urlID   string
	}{
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "reducing JSON test #1",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex1.com"},
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   "dc605989f530a3dfe9f7edacf1b3965b",
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "duplicate test #2",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex.com"},
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   "31aa70fc8589c52a763a2df36f304d28",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(strings.TrimSpace(utils.GetDBAddress())) != 0 {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				s := mock_repositories.NewMockURLRepository(ctrl)

				ctx := context.Background()

				if tt.wantErr {
					s.EXPECT().Save(ctx, tt.userID, tt.urlID, tt.request.URL).Return(fmt.Errorf("error"))
					u := New(s)
					_, err := u.ReduceURLToJSON(ctx, tt.userID, tt.request)
					assert.NotNil(t, err)
				} else {
					s.EXPECT().Save(ctx, tt.userID, tt.urlID, tt.request.URL).Return(nil)
					u := New(s)
					got, _ := u.ReduceURLToJSON(ctx, tt.userID, tt.request)
					assert.Equal(t, tt.want, got)
				}
			} else {
				u := New(tt.repo)
				got, _ := u.ReduceURLToJSON(nil, tt.userID, tt.request)
				if tt.wantErr {
					_, err := u.ReduceURLToJSON(nil, tt.userID, tt.request)
					assert.NotNil(t, err)
				} else {
					assert.Equalf(t, tt.want, got, "ReduceURLToJSON(%v)", tt.request)
				}
			}
		})
		testutils.AfterTest()
	}
}

func Test_urlServiceImpl_GetAllByUserID(t *testing.T) {
	utils.LoadEnvironments("../../../.env")

	tests := []struct {
		repo    repositories.URLRepository
		name    string
		url     string
		want    []entity.FullURL
		wantErr bool
		userID  string
	}{
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "get all test #1",
			want:    []entity.FullURL{{ShortURL: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b", OriginalURL: "yandex1.com"}},
			url:     "yandex1.com",
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
		{
			repo:    utils.CreateRepository(utils.GetFilePath(), utils.GetDBAddress()),
			name:    "not found test #2",
			url:     "yandex1.com",
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(strings.TrimSpace(utils.GetDBAddress())) != 0 {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				s := mock_repositories.NewMockURLRepository(ctrl)

				ctx := context.Background()

				if tt.wantErr {
					s.EXPECT().GetAllByUserID(ctx, tt.userID).Return(nil, fmt.Errorf("error"))
					u := New(s)
					_, err := u.GetAllByUserID(ctx, tt.userID)
					assert.NotNil(t, err)
				} else {
					s.EXPECT().GetAllByUserID(ctx, tt.userID).Return(tt.want, nil)
					u := New(s)
					got, _ := u.GetAllByUserID(ctx, tt.userID)
					assert.Equal(t, tt.want, got)
				}
			} else {
				u := New(tt.repo)
				u.ReduceAndSaveURL(nil, tt.userID, tt.url)

				got, _ := u.GetAllByUserID(nil, tt.userID)
				if tt.wantErr {
					wrongUserID := "dec27dda-6249-4f49-be71-4f56fc5ee541"
					_, err := u.GetAllByUserID(nil, wrongUserID)
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, tt.want, got)
				}
			}
		})
		testutils.AfterTest()
	}
}
