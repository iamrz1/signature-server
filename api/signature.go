package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"signature-server/data"
	_ "signature-server/docs"
	cerror "signature-server/error"
	cjson "signature-server/json"
	"signature-server/logger"
	"signature-server/middleware"
	"signature-server/model"

	"github.com/go-chi/chi"
	chiware "github.com/go-chi/chi/middleware"
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
	api.Use(logger.GenReqID)
	api.Use(chiware.Logger)
	api.Use(middleware.RequestLogger())
	api.Use(middleware.ResponseLogger())
}

func (api *signature) registerEndpoints() {
	api.Group(func(r chi.Router) {
		r.Get("/public_key", api.getPublicKey)
		r.Put("/transaction", api.createTransaction)
		r.Post("/signature", api.signTransactions)
		r.Post("/verify", api.verifySignature)
	})
}

// getPublicKey godoc
// @Summary Get public key corresponding to the daemon_key
// @Description Returns a JSON object containing the public key of the daemon key
// @Tags Signature
// @Produce  json
// @Success 200 {object} model.PublicKeyRes
// @Failure 400 {object} cjson.GenericErrorResponse
// @Failure 422 {object} cjson.GenericErrorResponse
// @Failure 500 {object} cjson.GenericErrorResponse
// @Router /public_key [get]
func (api *signature) getPublicKey(w http.ResponseWriter, r *http.Request) {
	cjson.ServeData(w, model.PublicKeyRes{PublicKey: api.sStore.GetPublicKey()})
}

// createTransaction godoc
// @Summary Create a transaction record
// @Description Takes a blob of data (arbitrary bytes) representing the transaction data in the form of a base64 string, and remembers it in memory. Returns a random, unique identifier for the transaction.
// @Tags Transaction
// @Accept  json
// @Produce  json
// @Param Body body model.CreateTransactionReq true "All fields are mandatory"
// @Success 200 {object} model.CreateTransactionRes
// @Failure 400 {object} cjson.GenericErrorResponse
// @Failure 422 {object} cjson.GenericErrorResponse
// @Failure 500 {object} cjson.GenericErrorResponse
// @Router /transaction [put]
func (api *signature) createTransaction(w http.ResponseWriter, r *http.Request) {
	body := &model.CreateTransactionReq{}
	if err := cjson.ParseBody(r, body); err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusBadRequest, "Unable to parse body", err))
		return
	}
	if err := body.Validate(); err != nil {
		logger.Errorf("Invalid request body %+v", err)
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
	cjson.ServeData(w, model.CreateTransactionRes{ID: id})
}

// signTransactions godoc
// @Summary Sign and return transactions
// @Description Takes a list of transaction identifiers, and builds a JSON array of strings containing the base64-encoded transaction blobs indicated by the given identifiers. It signs this array (serialised as JSON without any whitespace) using the daemon private key. Finally, it returns the array that was signed, as well as the signature as a base64 string.
// @Tags Signature
// @Accept  json
// @Produce  json
// @Param Body body model.SignTransactionReq true "All fields are mandatory"
// @Success 200 {object} model.SignTransactionRes
// @Failure 400 {object} cjson.GenericErrorResponse
// @Failure 422 {object} cjson.GenericErrorResponse
// @Failure 500 {object} cjson.GenericErrorResponse
// @Router /signature [post]
func (api *signature) signTransactions(w http.ResponseWriter, r *http.Request) {
	body := &model.SignTransactionReq{}
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

	messages := make([]string, 0)
	message, err := json.Marshal(txns)
	if err != nil {
		cjson.ServeError(w, cerror.NewAPIError(http.StatusInternalServerError, "Something went wrong", err))
		return
	}
	for _, txn := range txns {
		msg := base64.StdEncoding.EncodeToString(txn)
		messages = append(messages, msg)
	}

	cjson.ServeData(w, model.SignTransactionRes{Message: messages, Signature: api.sStore.SignData(message)})
}

// verifySignature godoc
// @Summary Verify a message using signature
// @Description Takes a message and the corresponding signature to verify it. Returns true or false as verification result.
// @Tags Signature
// @Accept  json
// @Produce  json
// @Param Body body model.VerifySignatureReq true "All fields are mandatory"
// @Success 200 {object} model.VerifySignatureRes
// @Failure 400 {object} cjson.VerifySignatureRes
// @Failure 422 {object} cjson.VerifySignatureRes
// @Failure 500 {object} cjson.VerifySignatureRes
// @Router /verify [post]
func (api *signature) verifySignature(w http.ResponseWriter, r *http.Request) {
	body := &model.VerifySignatureReq{}
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

	cjson.ServeData(w, model.VerifySignatureRes{Message: res})
}
