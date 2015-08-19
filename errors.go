package configo

import "errors"

var (
	KeyNotFoundError    = errors.New("The specified key was not found.")
	MalformedValueError = errors.New("The specified value could not be parsed.")
)
