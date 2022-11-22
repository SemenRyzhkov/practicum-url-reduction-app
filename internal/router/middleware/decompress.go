package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

func DecompressRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get(`Content-Encoding`) != `gzip` {
			next.ServeHTTP(w, r)
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		r.Body = io.NopCloser(gz)
		next.ServeHTTP(w, r)
		defer gz.Close()
	})
}
