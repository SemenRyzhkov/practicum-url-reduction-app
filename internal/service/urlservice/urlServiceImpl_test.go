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
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
			name:    "positive test #1",
			want:    "yandex.com",
			urlID:   "31aa70fc8589c52a763a2df36f304d28",
			wantErr: false,
		},
		{
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
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
				u.ReduceAndSaveURL(context.TODO(), tt.userID, tt.want)

				got, err := u.GetURLByID(context.TODO(), tt.urlID)
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
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
			name:    "positive test #1",
			saveURL: "yandex1.com",
			want:    "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b",
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   "dc605989f530a3dfe9f7edacf1b3965b",
		},
		{
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
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
				got, _ := u.ReduceAndSaveURL(context.TODO(), tt.userID, tt.saveURL)
				if tt.wantErr {
					_, err := u.ReduceAndSaveURL(context.TODO(), tt.userID, tt.saveURL)
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
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
			name:    "reducing JSON test #1",
			want:    entity.URLResponse{Result: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b"},
			request: entity.URLRequest{URL: "yandex1.com"},
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   "dc605989f530a3dfe9f7edacf1b3965b",
		},
		{
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
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
				got, _ := u.ReduceURLToJSON(context.TODO(), tt.userID, tt.request)
				if tt.wantErr {
					_, err := u.ReduceURLToJSON(context.TODO(), tt.userID, tt.request)
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
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
			name:    "get all test #1",
			want:    []entity.FullURL{{ShortURL: "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b", OriginalURL: "yandex1.com"}},
			url:     "yandex1.com",
			wantErr: false,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
		{
			repo:    utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
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
				u.ReduceAndSaveURL(context.TODO(), tt.userID, tt.url)

				got, _ := u.GetAllByUserID(context.TODO(), tt.userID)
				if tt.wantErr {
					wrongUserID := "dec27dda-6249-4f49-be71-4f56fc5ee541"
					_, err := u.GetAllByUserID(context.TODO(), wrongUserID)
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, tt.want, got)
				}
			}
		})
		testutils.AfterTest()
	}
}

func Test_urlServiceImpl_ReduceSeveralURL(t *testing.T) {
	utils.LoadEnvironments("../../../.env")

	tests := []struct {
		repo         repositories.URLRepository
		name         string
		request      []entity.URLWithIDRequest
		wantResponse []entity.URLWithIDResponse
		wantErr      bool
		userID       string
		urlID        []string
	}{
		{
			repo: utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
			name: "reduce several url test #1",
			request: []entity.URLWithIDRequest{
				{CorrelationID: "test1", OriginalURL: "yandex1.ru"},
				{CorrelationID: "test2", OriginalURL: "yandex2.ru"},
			},
			wantResponse: []entity.URLWithIDResponse{
				{
					CorrelationID: "test1",
					ShortURL:      "http://localhost:8080/b6ad61b613c33a6d62e6d14198e465b8",
				},
				{
					CorrelationID: "test2",
					ShortURL:      "http://localhost:8080/50754651b2f907807de0b789248f1f1b",
				},
			},
			wantErr: false,
			urlID:   []string{"b6ad61b613c33a6d62e6d14198e465b8", "50754651b2f907807de0b789248f1f1b"},
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
		},
		{
			repo: utils.CreateMemoryOrFileRepository(utils.GetFilePath()),
			name: "duplicate error test #2",
			request: []entity.URLWithIDRequest{
				{CorrelationID: "test1", OriginalURL: "yandex1.ru"},
				{CorrelationID: "test2", OriginalURL: "yandex1.ru"},
			},
			wantErr: true,
			userID:  "dec27dda-6249-4f49-be71-4f56fc5ee540",
			urlID:   []string{"b6ad61b613c33a6d62e6d14198e465b8"},
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
					s.EXPECT().Save(ctx, tt.userID, tt.urlID[0], tt.request[0].OriginalURL).Return(fmt.Errorf("error"))
					u := New(s)
					_, err := u.ReduceSeveralURL(ctx, tt.userID, tt.request)
					assert.NotNil(t, err)
				} else {
					s.EXPECT().Save(ctx, tt.userID, tt.urlID[0], tt.request[0].OriginalURL).Return(nil)
					s.EXPECT().Save(ctx, tt.userID, tt.urlID[1], tt.request[1].OriginalURL).Return(nil)
					u := New(s)
					got, _ := u.ReduceSeveralURL(ctx, tt.userID, tt.request)
					assert.Equal(t, tt.wantResponse, got)
				}
			} else {
				u := New(tt.repo)

				if tt.wantErr {
					_, err := u.ReduceSeveralURL(context.TODO(), tt.userID, tt.request)
					assert.NotNil(t, err)
				} else {
					got, _ := u.ReduceSeveralURL(context.TODO(), tt.userID, tt.request)
					assert.Equal(t, tt.wantResponse, got)
				}
			}
		})
		testutils.AfterTest()
	}
}
