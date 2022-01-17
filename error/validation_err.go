package cerror

import (
	"encoding/json"
	"fmt"
)

var (
	InvalidInputErr = fmt.Errorf("input data is not valid")
	NotFoundErr     = fmt.Errorf("could not find data")
)

// ValidationError ...
type ValidationError map[string]string

func (err ValidationError) Error() string {
	buf, _ := json.Marshal(err)
	return string(buf)
}

// Add ...
func (err ValidationError) Add(key, msg string) {
	err[key] = msg
}
