package configo

type MultiSource []Source

func (this MultiSource) Initialize() {
	for _, source := range this {
		source.Initialize()
	}
}

func (this MultiSource) Strings(key string) (result []string, err error) {
	for _, source := range this {
		result, err = source.Strings(key)
		if err == nil {
			return result, err
		}
	}
	return nil, KeyNotFoundError
}
