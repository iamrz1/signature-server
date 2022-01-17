package cjson

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"net/http"
)

// ParseBody ...
func ParseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ParseQuery ...
func ParseQuery(r *http.Request, v interface{}) error {
	return schema.NewDecoder().Decode(v, r.URL.Query())
}
