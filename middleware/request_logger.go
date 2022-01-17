package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"signature-server/util"
)

// RequestLogger ...
func RequestLogger(enable bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enable {
				util.Infof("Req: %v %v", r.Method, r.URL)
				if body, err := readIntact(r); err == nil {
					util.Infof("Body: %v", string(body))
				} else {
					util.Error(err.Error())
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func readIntact(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	body, err := ioutil.ReadAll(tee)
	r.Body = ioutil.NopCloser(&buf)
	return body, err
}
