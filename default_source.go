package newton

import (
	"net/url"
	"strconv"
	"time"
)

type DefaultSource struct {
	settings map[string][]string
}

func NewDefaultSource() *DefaultSource {
	return &DefaultSource{
		settings: make(map[string][]string),
	}
}

func (this *DefaultSource) Add(key string, values ...interface{}) *DefaultSource {
	contents := this.settings[key]

	for _, value := range values {
		contents = append(contents, convertToString(value))
	}

	this.settings[key] = contents
	return this
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

func (this *DefaultSource) Strings(key string) ([]string, error) {
	values, found := this.settings[key]
	if !found {
		return nil, KeyNotFoundError
	}

	return values, nil
}

func (this *DefaultSource) Initialize() {}
