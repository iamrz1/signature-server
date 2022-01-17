package middleware

import (
	"net/http"
	"net/http/httptest"
	"signature-server/util"
	"strings"
)

// ResponseLogger ...
func ResponseLogger(enable bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enable {
				res := httptest.NewRecorder()
				next.ServeHTTP(res, r)
				util.Infof("Resp: %v", strings.TrimSpace(res.Body.String()))

				for k, v := range res.HeaderMap {
					w.Header()[k] = v
				}
				w.WriteHeader(res.Code)
				res.Body.WriteTo(w)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
