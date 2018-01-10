package configo

import (
	"net/url"
	"strconv"
	"time"
)

// DefaultSource is allows registration of specified default values of various types.
type DefaultSource struct {
	settings map[string][]string
}

// NewDefaultSource initializes a new DefaultSource.
func NewDefaultSource(pairs ...DefaultPair) *DefaultSource {
	source := &DefaultSource{settings: make(map[string][]string)}
	for _, config := range pairs {
		config(source)
	}
	return source
}

type DefaultPair func(*DefaultSource)

// Default registers the provided values (which will be converted to strings) to the given key.
// It does NOT overwrite existing values, it appends.
func Default(key string, values ...interface{}) DefaultPair {
	return func(source *DefaultSource) {
		contents := source.settings[key]

		for _, value := range values {
			contents = append(contents, convertToString(value))
		}

		source.settings[key] = contents
	}
}

func convertToString(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return typed
	case bool:
		return strconv.FormatBool(typed)
	case int:
		return strconv.FormatInt(int64(typed), 10)
	case int8:
		return strconv.FormatInt(int64(typed), 10)
	case int16:
		return strconv.FormatInt(int64(typed), 10)
	case int32:
		return strconv.FormatInt(int64(typed), 10)
	case int64:
		return strconv.FormatInt(typed, 10)
	case uint:
		return strconv.FormatUint(uint64(typed), 10)
	case uint8:
		return strconv.FormatUint(uint64(typed), 10)
	case uint16:
		return strconv.FormatUint(uint64(typed), 10)
	case uint32:
		return strconv.FormatUint(uint64(typed), 10)
	case uint64:
		return strconv.FormatUint(typed, 10)
	case float32:
		return strconv.FormatFloat(float64(typed), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case *url.URL:
		return typed.String()
	case url.URL:
		return typed.String()
	case time.Duration:
		return typed.String()
	case time.Time:
		return typed.String()
	default:
		return ""
	}
}

// Strings returns all values associated with the given key, or KeyNotFoundError.
func (this *DefaultSource) Strings(key string) ([]string, error) {
	values, found := this.settings[key]
	if !found {
		return nil, KeyNotFoundError
	}

	return values, nil
}

func (this *DefaultSource) Initialize() {}
