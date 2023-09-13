package cjson

import (
	"encoding/json"
	"net/http"
)

// ParseBody ...
func ParseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
