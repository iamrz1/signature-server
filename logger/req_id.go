package logger

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"github.com/tylerb/gls"
)

// ReqIDTag holds tag for request id
const ReqIDTag = "RequestID"

type ridHook struct{}

// Run gets RequestID from middleware and sets it in the event context
func (h ridHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if id := gls.Get(ReqIDTag); id != nil && level != zerolog.NoLevel {
		e.Str(ReqIDTag, fmt.Sprintf("%v", id))
	}
}

// GenReqID generates a new request id for each request and stores in the gls. Storage is cleaned up after request is server.
func GenReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gls.Set(ReqIDTag, uuid.NewV4().String())
		defer gls.Cleanup()
		next.ServeHTTP(w, r)
	})
}
