package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

func Test_urlServiceImpl_GetUrlById(t *testing.T) {
	tests := []struct {
		repo    repositories.UrlRepository
		name    string
		urlId   string
		want    string
		wantErr bool
	}{
		{
			repo:    repositories.NewUrlRepository(),
			name:    "positive test #1",
			want:    "yandex.com",
			urlId:   "31aa70fc8589c52a763a2df36f304d28",
			wantErr: false,
		},
		{
			repo:    repositories.NewUrlRepository(),
			name:    "not found test #2",
			want:    "yandex.com",
			urlId:   "31aa70fc8589c52a763a2df36f304d29",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlService(tt.repo)
			u.ReduceAndSaveUrl(tt.want)

			got, err := u.GetUrlById(tt.urlId)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, got)

			}
		})
	}
}

func Test_urlServiceImpl_ReduceAndSaveUrl(t *testing.T) {
	tests := []struct {
		repo    repositories.UrlRepository
		name    string
		saveUrl string
		want    string
		wantErr bool
	}{
		{
			repo:    repositories.NewUrlRepository(),
			name:    "positive test #1",
			saveUrl: "yandex1.com",
			want:    "http://localhost:8080/dc605989f530a3dfe9f7edacf1b3965b",
			wantErr: false,
		},
		{
			repo:    repositories.NewUrlRepository(),
			name:    "duplicate test #2",
			saveUrl: "yandex.com",
			want:    "http://localhost:8080/XVlBz",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlService(tt.repo)
			got, _ := u.ReduceAndSaveUrl(tt.saveUrl)
			if tt.wantErr {
				_, err := u.ReduceAndSaveUrl(tt.saveUrl)
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
