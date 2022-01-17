package util

import (
	"net/http"
	"testing"
)

func testHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		Info("Hello World!")
	}
	return fn
}

func TestGatekeeper(t *testing.T) {
	svr := GenReqID(testHandler())
	svr.ServeHTTP(nil, nil)
}
