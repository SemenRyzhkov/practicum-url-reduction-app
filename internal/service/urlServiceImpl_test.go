package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_urlServiceImpl_GetUrlById(t *testing.T) {
	tests := []struct {
		name    string
		urlId   string
		want    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			want:    "yandex.com",
			urlId:   "XVlBz",
			wantErr: false,
		},
		{
			name:    "not found test #2",
			want:    "yandex.com",
			urlId:   "1111",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlService()
			u.ReduceAndSaveUrl(tt.want)

			got, err := u.GetUrlById(tt.urlId)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, got, tt.want)

			}
		})
	}
}

func Test_urlServiceImpl_ReduceAndSaveUrl(t *testing.T) {
	tests := []struct {
		name    string
		saveUrl string
		want    string
		wantErr bool
	}{
		{
			name:    "positive test #1",
			saveUrl: "yandex1.com",
			want:    "http://localhost:8080/MRAjW",
			wantErr: false,
		},
		{
			name:    "duplicate test #2",
			saveUrl: "yandex.com",
			want:    "http://localhost:8080/XVlBz",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlService()

			got, err := u.ReduceAndSaveUrl(tt.saveUrl)
			fmt.Println(err)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
