package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		start := time.Now()

		defer func() {
			panicVal := recover()
			if panicVal != nil {
				lrw.statusCode = http.StatusInternalServerError
				panic(panicVal)
			}

			duration_sec := time.Since(start).Seconds()
			log.Info().Ctx(r.Context()).
				Str("method", r.Method).
				Stringer("path", r.URL).
				Int("status_code", lrw.statusCode).
				Str("host", r.Host).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Str("referer", r.Referer()).
				Float64("duration_sec", duration_sec).
				Msg("Access log")
		}()

		next.ServeHTTP(lrw, r)
	})
}
