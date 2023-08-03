package middleware

import (
	"app/util"
	"errors"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func Recover(next http.Handler) http.Handler {
	logger := zap.S()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						seErr := strings.ToLower(se.Error())
						if strings.Contains(seErr, "broken pipe") ||
							strings.Contains(seErr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(r, false)
				if brokenPipe {
					logger.Error(r.URL.Path, zap.Any("error", err))
					util.Reject(w, -1, errors.New("broken pipe"))
					return
				}
				logger.Errorf("req %s recover from error %v at %v\n", string(httpRequest), err, time.Now())
				util.Reject(w, -1, err)

			}
		}()
		next.ServeHTTP(w, r)
	})
}
