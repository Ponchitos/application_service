package middleware

import (
	"github.com/Ponchitos/application_service/server/tools/logger"
	"net/http"
	"time"
)

func LoggerRequest(lgr logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lgr.Infof("started %s %s remote_addr:%v request_id:%v", r.Method, r.RequestURI, r.RemoteAddr, r.Context().Value(CtxKeyRequestID))

			start := time.Now()
			next.ServeHTTP(w, r)

			lgr.Infof(
				"completed with in %v. remote_addr:%v request_id:%v",
				time.Since(start),
				r.RemoteAddr,
				r.Context().Value(CtxKeyRequestID),
			)
		})
	}
}
