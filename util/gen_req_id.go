package util

import (
	uuid "github.com/satori/go.uuid"
	"github.com/tylerb/gls"
	"net/http"
)

// ReqIDTag holds tag for request id
const ReqIDTag = "request_id"

// GenReqID generates a new request id for each requests and stores in the gls. Storage is cleaned up after request is server.
func GenReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gls.Set(ReqIDTag, uuid.NewV4().String())
		defer gls.Cleanup()
		next.ServeHTTP(w, r)
	})
}
