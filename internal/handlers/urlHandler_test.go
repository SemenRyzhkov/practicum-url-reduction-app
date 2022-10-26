package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"testing"
)

func TestGetUrlById(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
		params  httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetUrlById(tt.args.writer, tt.args.request, tt.args.params)
		})
	}
}

func TestReduceUrl(t *testing.T) {
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
		in2     httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReduceUrl(tt.args.writer, tt.args.request, tt.args.in2)
		})
	}
}
