package newton

import (
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"
)

type ConditionalSource struct {
	condition func() bool
	settings  map[string][]string
}

func NewDevelopmentSource() *ConditionalSource {
	return NewConditionalSource(func() bool {
		hostname, _ := os.Hostname()
		return runtime.GOOS == "darwin" || hostname == "vagrant"
	})
}

func NewConditionalSource(condition func() bool) *ConditionalSource {
	return &ConditionalSource{
		condition: condition,
		settings:  make(map[string][]string),
	}
}

func (this *ConditionalSource) Add(key string, values ...interface{}) {
	contents := this.settings[key]

	for _, value := range values {
		contents = append(contents, convertToString(value))
	}

	this.settings[key] = contents
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

func (this *ConditionalSource) Strings(key string) ([]string, error) {
	if !this.condition() {
		return nil, KeyNotFoundError
	}

	values, found := this.settings[key]
	if !found {
		return nil, KeyNotFoundError
	}

	return values, nil
}
