package memory

import (
	"encoding/base64"
	"signature-server/data"
	cerror "signature-server/error"
	"sync"
)

type transactionStore struct {
	id   int64
	mp   map[int64][]byte
	lock sync.Mutex
}

// NewTransactionStore ...
func NewTransactionStore() data.TransactionStore {
	return &transactionStore{
		id:   0,
		mp:   map[int64][]byte{},
		lock: sync.Mutex{},
	}
}

func (s *transactionStore) InsertOne(data string) (int64, error) {
	blob, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return 0, cerror.InvalidInputErr
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.id++
	s.mp[s.id] = blob

	return s.id, nil
}

func (s *transactionStore) FindOne(id int64) (string, error) {
	blob, ok := s.mp[id]
	if !ok {
		return "", cerror.NotFoundErr
	}
	return base64.StdEncoding.EncodeToString(blob), nil
}

func (s *transactionStore) FindOneBlob(id int64) ([]byte, error) {
	blob, ok := s.mp[id]
	if !ok {
		return nil, cerror.NotFoundErr
	}
	return blob, nil
}

func (s *transactionStore) FindMany(ids []int64) ([]string, error) {
	res := make([]string, 0)
	var notFound bool
	for _, id := range ids {
		blob, ok := s.mp[id]
		if !ok {
			notFound = true
			break
		}
		res = append(res, base64.StdEncoding.EncodeToString(blob))
	}

	if notFound {
		return nil, cerror.NotFoundErr
	}

	return res, nil
}

func (s *transactionStore) FindManyBlobs(ids []int64) ([][]byte, error) {
	res := make([][]byte, 0)
	var notFound bool
	for _, id := range ids {
		blob, ok := s.mp[id]
		if !ok {
			notFound = true
			break
		}
		res = append(res, blob)
	}

	if notFound {
		return nil, cerror.NotFoundErr
	}

	return res, nil
}
