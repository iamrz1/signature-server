package cerror

import (
	"fmt"
)

// APIError ...
type APIError struct {
	Status  int
	Message string
	tags    []string
	source  error
}

// NewAPIError ...
func NewAPIError(status int, msg string, src error, tags ...string) *APIError {
	err := &APIError{
		Status:  status,
		Message: msg,
		tags:    tags,
		source:  src,
	}
	return err
}

func (err *APIError) Error() string {
	return fmt.Sprintf("%v: %v", err.Message, err.source)
}
