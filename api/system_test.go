package api_test

import (
	"net/http"
	"net/http/httptest"
	"signature-server/api"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

func TestSystemCheck(t *testing.T) {
	h := api.NewSystemHandler()
	t.Run("system check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/check", nil)
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		jsonassert.New(t).Assertf(res.Body.String(), `{"system_status":"ok"}`)
	})
}

func TestSystemPanic(t *testing.T) {
	h := api.NewSystemHandler()
	t.Run("system panic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		jsonassert.New(t).Assertf(res.Body.String(), `{"message":"system panic"}`)
	})
}
