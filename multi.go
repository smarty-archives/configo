package configo

type multiSource struct {
	sources []Source
}

func MultiSource(sources ...Source) *multiSource {
	return &multiSource{sources: sources}
}

func (this *multiSource) Initialize() {
	for _, source := range this.sources {
		source.Initialize()
	}
}

func (this *multiSource) Strings(key string) (result []string, err error) {
	for _, source := range this.sources {
		result, err = source.Strings(key)
		if err == nil {
			return result, err
		}
	}
	return nil, KeyNotFoundError
}
