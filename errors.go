package configo

import "errors"

var (
	ErrKeyNotFound    = errors.New("the specified key was not found")
	ErrMalformedValue = errors.New("the specified value could not be parsed")
)
