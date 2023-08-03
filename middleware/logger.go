package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	logger := zap.S()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		query := r.URL.RawQuery
		next.ServeHTTP(w, r)
		end := time.Now()
		latency := end.Sub(start)
		logger.Infof("%s %s?%s latency %v at %v", r.Method, path, query, latency, end.Format("2006-01-02 15:04:05"))
	})
}
