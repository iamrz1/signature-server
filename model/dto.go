package model

type PublicKeyRes struct {
	PublicKey string `json:"public_key"`
}

type CreateTransactionReq struct {
	TxnData string `json:"txn"`
}

type CreateTransactionRes struct {
	ID int64 `json:"id"`
}

type SignTransactionReq struct {
	IDs []int64 `json:"ids"`
}

type SignTransactionRes struct {
	Message   []string `json:"message"`
	Signature string   `json:"signature"`
}

type VerifySignatureReq struct {
	Message   interface{} `json:"message"`
	Signature string      `json:"signature"`
}

type VerifySignatureRes struct {
	Message bool `json:"message"`
}
