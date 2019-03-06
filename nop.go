package configo

import "reflect"

func FirstOrNop(sources ...Source) Source {
	for _, source := range sources {
		if !isNil(source) {
			return source
		}
	}
	return new(nopSource)
}

func isNil(source Source) bool {
	value := reflect.ValueOf(source)
	return !value.IsValid() || value.IsNil()
}

type nopSource struct{}

func (*nopSource) Initialize() {}

func (*nopSource) Strings(key string) ([]string, error) { return nil, KeyNotFoundError }
