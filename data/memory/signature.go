package memory

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"signature-server/data"
)

type signatureStore struct {
	DaemonKey       string
	PublicKeyString string
	PrivateKey      []byte
	Seed            []byte
	PublicKey       []byte
}

func NewSignatureStore(daemonKey string) (data.SignatureStore, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(daemonKey)
	if err != nil {
		return nil, fmt.Errorf("invalid daemon key")
	}
	if len(keyBytes) < ed25519.SeedSize {
		return nil, fmt.Errorf("daemon key must be atlest 32 bytes long")
	}
	s := &signatureStore{
		DaemonKey: daemonKey,
		Seed:      keyBytes[:ed25519.SeedSize],
	}

	s.PrivateKey = ed25519.NewKeyFromSeed(s.Seed)
	s.PublicKey = make([]byte, ed25519.PublicKeySize)
	n := copy(s.PublicKey, s.PrivateKey[ed25519.SeedSize:])
	if n != 32 {
		return nil, fmt.Errorf("could not generate 32 bytes long public key")
	}

	s.PublicKeyString = base64.StdEncoding.EncodeToString(s.PublicKey)

	return s, nil
}

func (s *signatureStore) GetPublicKey() string {
	return s.PublicKeyString
}

func (s *signatureStore) SignData(data []byte) string {
	signature := ed25519.Sign(s.PrivateKey, data)

	return base64.StdEncoding.EncodeToString(signature)
}

func (s *signatureStore) VerifySignature(data, signature []byte) bool {
	return ed25519.Verify(s.PublicKey, data, signature)
}
