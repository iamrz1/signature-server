package data

// TransactionStore ...
type TransactionStore interface {
	InsertOne(data string) (int64, error)
	FindOne(id int64) (string, error)
	FindOneBlob(id int64) ([]byte, error)
	FindMany(ids []int64) ([]string, error)
	FindManyBlobs(ids []int64) ([][]byte, error)
}
