package cjson

import (
	"encoding/json"
	"net/http"
	cerror "signature-server/error"
)

// Object ...
type Object map[string]interface{}

// Add ...
func (o Object) Add(key string, msg interface{}) {
	o[key] = msg
}

type response struct {
	status  int
	err     bool
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (res *response) serveJSON(w http.ResponseWriter) {
	if res.status == 0 {
		res.status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.status)
	if res.err {
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(res.Data); err != nil {
		panic(err)
	}
}

// ServeJSON ...
func ServeJSON(w http.ResponseWriter, data interface{}, status int) {
	res := response{
		status: status,
		Data:   data,
	}
	res.serveJSON(w)
}

// ServeData ...
func ServeData(w http.ResponseWriter, data interface{}) {
	ServeJSON(w, data, http.StatusOK)
}

// ServeError ...
func ServeError(w http.ResponseWriter, err *cerror.APIError) {
	res := response{
		err:     true,
		status:  err.Status,
		Message: err.Message,
	}
	res.serveJSON(w)
}

type GenericErrorResponse struct {
	Message string
}
