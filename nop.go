package configo

func FirstOrNop(sources ...Source) Source {
	for _, source := range sources {
		if source != nil {
			return source
		}
	}
	return new(nopSource)
}

type nopSource struct{}

func (*nopSource) Initialize() {}

func (*nopSource) Strings(key string) ([]string, error) { return nil, KeyNotFoundError }


