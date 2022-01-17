package data

// SignatureStore ...
type SignatureStore interface {
	GetPublicKey() string
	SignData(data []byte) string
	VerifySignature(data, signature []byte) bool
}
