package configo

import "errors"

var (
	KeyNotFoundError    = errors.New("the specified key was not found")
	MalformedValueError = errors.New("the specified value could not be parsed")
)
