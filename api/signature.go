package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
	"net/http"
	"signature-server/data"
	cerror "signature-server/error"
	cjson "signature-server/json"
	"signature-server/middleware"
	"signature-server/util"
)

type signature struct {
	chi.Router
	sStore data.SignatureStore
	tStore data.TransactionStore
}

// NewSignatureHandler ...
func NewSignatureHandler(sStore data.SignatureStore, tStore data.TransactionStore) http.Handler {
	h := &signature{
		chi.NewRouter(),
		sStore,
		tStore,
	}
	h.registerMiddleware()
	h.registerEndpoints()
	return h
}

func (api *signature) registerMiddleware() {
	api.Use(util.GenReqID)
	api.Use(chiware.Logger)
	api.Use(middleware.RequestLogger(true))
	api.Use(middleware.ResponseLogger(true))
}

func (api *signature) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Get("/public_key", api.getPublicKey)
		r.Put("/transaction", api.createTransaction)
		r.Post("/signature", api.signTransactions)
		r.Post("/verify", api.verifySignature)
	})
}

func (api *signature) getPublicKey(w http.ResponseWriter, r *http.Request) {
	cjson.ServeData(w, cjson.Object{"public_key": api.sStore.GetPublicKey()})
}

func (api *signature) createTransaction(w http.ResponseWriter, r *http.Request) {
	body := &createTransactionRequestBody{}
	if err := cjson.ParseBody(r, body); err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusBadRequest, "Unable to parse body", err))
		return
	}
	if err := body.Validate(); err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusUnprocessableEntity, "Invalid request body", err))
		return
	}

	id, err := api.tStore.InsertOne(body.TxnData)
	if err != nil {
		if err == cerror.InvalidInputErr {
			cjson.ServeError(w, cerror.NewAPIError(http.StatusUnprocessableEntity, "Invalid request body", err))
			return
		}
		cjson.ServeError(w, cerror.NewAPIError(http.StatusInternalServerError, "Something went wrong", err))
		return
	}
	cjson.ServeData(w, cjson.Object{"id": id})
}

type createTransactionRequestBody struct {
	TxnData string `json:"txn"`
}

func (b *createTransactionRequestBody) Validate() error {
	err := cerror.ValidationError{}
	if b.TxnData == "" {
		err.Add("txn", "required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (api *signature) signTransactions(w http.ResponseWriter, r *http.Request) {
	body := &signTransactionRequestBody{}
	if err := cjson.ParseBody(r, body); err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusBadRequest, "Unable to parse body", err))
		return
	}
	txns, err := api.tStore.FindManyBlobs(body.IDs)
	if err != nil {
		if err == cerror.NotFoundErr {
			cjson.ServeError(w, cerror.NewAPIError(http.StatusUnprocessableEntity, "One or more transaction ID don't exist", err))
			return
		}
		cjson.ServeError(w, cerror.NewAPIError(http.StatusInternalServerError, "Something went wrong", err))
		return
	}

	message, err := json.Marshal(txns)
	if err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusInternalServerError, "Something went wrong", err))
		return
	}

	cjson.ServeData(w, cjson.Object{"message": txns, "signature": api.sStore.SignData(message)})
}

type signTransactionRequestBody struct {
	IDs []int64 `json:"ids"`
}

func (api *signature) verifySignature(w http.ResponseWriter, r *http.Request) {
	body := &verifySignatureRequestBody{}
	if err := cjson.ParseBody(r, body); err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusBadRequest, "Unable to parse body", err))
		return
	}
	if err := body.Validate(); err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusUnprocessableEntity, "Invalid request body", err))
		return
	}

	message, err := json.Marshal(body.Message)
	if err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusUnprocessableEntity, "Invalid message", err))
		return
	}

	sig, err := base64.StdEncoding.DecodeString(body.Signature)
	if err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusUnprocessableEntity, "Invalid signature", err))
		return
	}

	res := api.sStore.VerifySignature(message, sig)

	cjson.ServeData(w, cjson.Object{"message": res})
}

type verifySignatureRequestBody struct {
	Message   interface{} `json:"message"`
	Signature string      `json:"signature"`
}

func (b *verifySignatureRequestBody) Validate() error {
	err := cerror.ValidationError{}
	if b.Message == nil {
		err.Add("message", "required")
	}

	if b.Signature == "" {
		err.Add("signature", "required")
	}

	if len(err) != 0 {
		return err
	}

	return nil
}
