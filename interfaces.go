package configo

// Source defines the methods required by a Reader.
// The Strings method returns all values associated with the given key with an error
// if the key does not exist.
type Source interface {
	Initialize()
	Strings(key string) ([]string, error)
}
