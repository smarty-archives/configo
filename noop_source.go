package configo

type NoopSource struct{}

func (this NoopSource) Initialize() {
}
func (this NoopSource) Strings(string) ([]string, error) {
	return nil, KeyNotFoundError
}
