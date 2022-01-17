package api_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"signature-server/data/memory"

	"signature-server/api"
)

const (
	TestDaemonKey = "G7bxZbii5tgq1X8rgdOmy/zSwaOSKtHleSSw41YhJ2aeErC62pZtkQRBEA4xNqqTd+VouCf5/PPvTgW9seqyBg=="
)

func TestGetPublicKey(t *testing.T) {
	sStore, err := memory.NewSignatureStore(TestDaemonKey)
	if err != nil {
		t.Fatal(err)
	}
	tStore := memory.NewTransactionStore()
	h := api.NewSignatureHandler(sStore, tStore)
	testData := []struct {
		des  string
		code int
		res  string
	}{
		{
			des:  "get public key",
			code: http.StatusOK,
			res:  `{"public_key": "nhKwutqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY="}`,
		},
	}
	for _, td := range testData {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/public_key", nil)
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.res)
		})
	}
}

func TestInsertTransaction(t *testing.T) {
	sStore, err := memory.NewSignatureStore(TestDaemonKey)
	if err != nil {
		t.Fatal(err)
	}
	tStore := memory.NewTransactionStore()
	h := api.NewSignatureHandler(sStore, tStore)
	testData := []struct {
		des  string
		code int
		txn  string
		res  string
	}{
		{
			des:  "first transaction",
			code: http.StatusOK,
			txn:  "nhKwutqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY=",
			res:  `{"id": 1}`,
		},
		{
			des:  "faulty transaction data",
			code: http.StatusUnprocessableEntity,
			txn:  "nhKwutqWbZEEQRAOMTdaqk3fflaLgn+fzz704FvbHqsgY=",
			res:  `{"message": "Invalid request body"}`,
		},
		{
			des:  "second transaction",
			code: http.StatusOK,
			txn:  "/+ABAgM=",
			res:  `{"id": 2}`,
		},
	}
	for _, td := range testData {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/transaction", bytes.NewReader([]byte(fmt.Sprintf(`{"txn":"%s"}`, td.txn))))
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.res)
		})
	}
}

func TestSignature(t *testing.T) {
	sStore, err := memory.NewSignatureStore(TestDaemonKey)
	if err != nil {
		t.Fatal(err)
	}
	tStore := memory.NewTransactionStore()
	h := api.NewSignatureHandler(sStore, tStore)
	txns := []string{"nhKwutqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY=", "/+ABAgM="}
	for _, tx := range txns {
		req := httptest.NewRequest(http.MethodPut, "/transaction", bytes.NewReader([]byte(fmt.Sprintf(`{"txn":"%s"}`, tx))))
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)
	}
	testData := []struct {
		des  string
		code int
		body string
		res  string
	}{
		{
			des:  "correct signature request",
			code: http.StatusOK,
			body: `{"ids": [1]}`,
			res:  `{"message": ["nhKwutqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY="],"signature": "fCoelFB6mHNuYrCX+nT+ZgNeeqd182KTwiRpKBsJQRYpwIU/Db9uMbDiBYMl3XiadE7/vagUcVFt79V//3ykBQ=="}`,
		},
		{
			des:  "non-existing data signature",
			code: http.StatusUnprocessableEntity,
			body: `{"ids": [3]}`,
			res:  `{"message": "One or more transaction ID don't exist"}`,
		},
		{
			des:  "faulty request body",
			code: http.StatusBadRequest,
			body: `{"ids": 1}`,
			res:  `{"message": "Unable to parse body"}`,
		},
	}
	for _, td := range testData {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/signature", bytes.NewReader([]byte(td.body)))
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.res)
		})
	}
}

func TestVerifySignature(t *testing.T) {
	sStore, err := memory.NewSignatureStore(TestDaemonKey)
	if err != nil {
		t.Fatal(err)
	}
	tStore := memory.NewTransactionStore()
	h := api.NewSignatureHandler(sStore, tStore)
	txns := []string{"nhKwutqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY=", "/+ABAgM="}
	for _, tx := range txns {
		req := httptest.NewRequest(http.MethodPut, "/transaction", bytes.NewReader([]byte(fmt.Sprintf(`{"txn":"%s"}`, tx))))
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)
	}
	testData := []struct {
		des  string
		code int
		body string
		res  string
	}{
		{
			des:  "correct signature request",
			code: http.StatusOK,
			body: `{"message": ["nhKwutqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY="],"signature": "fCoelFB6mHNuYrCX+nT+ZgNeeqd182KTwiRpKBsJQRYpwIU/Db9uMbDiBYMl3XiadE7/vagUcVFt79V//3ykBQ=="}`,
			res:  `{"message": true}`,
		},
		{
			des:  "faulty signature in request",
			code: http.StatusUnprocessableEntity,
			body: `{"message": ["nhKwuteqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY="],"signature": "fCoelFB6fmHNuYrCX+nT+ZgNeeqd182KTwiRpKBsJQRYpwIU/Db9uMbDiBYMl3XiadE7/vagUcVFt79V//3ykBQ=="}`,
			res:  `{"message": "Invalid signature"}`,
		},
		{
			des:  "missing field in request",
			code: http.StatusUnprocessableEntity,
			body: `{"message": ["nhKwuteqWbZEEQRAOMTaqk3flaLgn+fzz704FvbHqsgY="]}`,
			res:  `{"message": "Invalid request body"}`,
		},
	}
	for _, td := range testData {
		t.Run(td.des, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewReader([]byte(td.body)))
			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)
			assert.Equal(t, td.code, res.Code)
			jsonassert.New(t).Assertf(res.Body.String(), td.res)
		})
	}
}
