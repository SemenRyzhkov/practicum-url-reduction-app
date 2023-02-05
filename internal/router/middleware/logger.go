package middleware

import (
	stdlog "log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-kit/log"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// Status статус
func (rw *responseWriter) Status() int {
	return rw.status
}

// WriteHeader врайтер хидера
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

}

func getLogger() log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)
	return logger
}

// LoggingMiddleware логирование для хэндлеров
func LoggingMiddleware(next http.Handler) http.Handler {
	loggerrr := getLogger()
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				loggerrr.Log(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		loggerrr.Log(
			"status", wrapped.status,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"duration", time.Since(start),
		)
	}

	return http.HandlerFunc(fn)

}
