package configo

import (
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Reader retrieves values from the provided sources, handling conversions
// to the type identified by the method being called (Strings, Ints, etc...).
type Reader struct {
	sources []Source
	aliases map[string][]string
	fatal   func(string, error)
}

// NewReader initializes a new reader using the provided sources. It calls each
// non-nil source's Initialize() method.
func NewReader(sources ...Source) *Reader {
	return &Reader{
		sources: initialize(sources),
		aliases: make(map[string][]string),
		fatal: func(key string, err error) {
			log.Fatalf("[%s] %s\n", key, err)
		},
	}
}
func initialize(sources []Source) (filtered []Source) {
	for _, source := range sources {
		if source == nil {
			continue
		}

		value := reflect.ValueOf(source)
		if value.Type().Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		source.Initialize()
		filtered = append(filtered, source)
	}
	return filtered
}

func (this *Reader) RegisterAlias(key, alias string) {
	this.aliases[alias] = append(this.aliases[alias], key)
	this.aliases[key] = append(this.aliases[key], alias)
}

// Strings returns all values associated with the given key or nil
// if the key does not exist.
func (this *Reader) Strings(key string) []string {
	value, _ := this.StringsError(key)
	return value
}

// StringsError returns all values associated with the given key with an error
// if the key does not exist. It does so by searching it sources, in the order
// they were provided, and returns the first non-error result or ErrKeyNotFound.
func (this *Reader) StringsError(key string) ([]string, error) {
	for _, alias := range this.resolvePossibleKeys(key) {
		if values, err := this.stringsError(alias); err == nil {
			return values, nil
		}
	}

	return nil, ErrKeyNotFound
}
func (this *Reader) stringsError(key string) ([]string, error) {
	for _, source := range this.sources {
		if value, err := source.Strings(key); err == nil {
			if err == nil && len(value) > 0 && strings.HasPrefix(value[0], "env:") {
				key = value[0] // if an EnvironmentSource is still to be inspected, it will remove the 'env:' prefix and do the lookup.
				continue
			}
			return value, nil
		}
	}

	return nil, ErrKeyNotFound
}

func (this *Reader) resolvePossibleKeys(key string) []string {
	return append([]string{key}, this.aliases[key]...)
}

// StringsPanic returns all values associated with the given key or panics
// if the key does not exist.
func (this *Reader) StringsPanic(key string) []string {
	if value, err := this.StringsError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// StringsDefault returns all values associated with the given key or the provided defaults
// if the key does not exist.
func (this *Reader) StringsDefault(key string, Default []string) []string {
	if value, err := this.StringsError(key); err != nil {
		return Default
	} else {
		return value
	}
}

// StringsFatal returns all values associated with the given key or log.Fatal()
// if the key does not exist.
func (this *Reader) StringsFatal(key string) []string {
	if value, err := this.StringsError(key); err != nil {
		this.fatal(key, err)
		return nil
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// String returns the first value associated with the given key or an empty string
// if the key does not exist.
func (this *Reader) String(key string) string {
	if value, err := this.StringError(key); err != nil {
		return ""
	} else {
		return value
	}
}

// StringError returns the first value associated with the given key with an an error
// if the key does not exist.
func (this *Reader) StringError(key string) (string, error) {
	if value, err := this.StringsError(key); err != nil {
		return "", err
	} else if len(value) == 0 {
		return "", nil
	} else {
		return value[0], nil
	}
}

// StringPanic returns the first value associated with the given key or panics
// if the key does not exist.
func (this *Reader) StringPanic(key string) string {
	if value, err := this.StringError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// StringDefault returns the first value associated with the given key or the provided default
// if the key does not exist.
func (this *Reader) StringDefault(key string, Default string) string {
	if value, err := this.StringError(key); err != nil {
		return Default
	} else {
		return value
	}
}

// StringFatal returns the first value associated with the given key or calls log.Fatal()
// if the key does not exist.
func (this *Reader) StringFatal(key string) string {
	if value, err := this.StringError(key); err != nil {
		this.fatal(key, err)
		return ""
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// Ints returns all integer values associated with the given key or returns 0
// if the key does not exist.
func (this *Reader) Ints(key string) []int {
	value, _ := this.IntsError(key)
	return value
}

// IntsError returns all integer values associated with the given key with an error
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntsError(key string) ([]int, error) {
	raw, err := this.StringsError(key)
	if err != nil {
		return nil, err
	}

	ints := make([]int, len(raw))
	for i, r := range raw {
		ints[i], err = strconv.Atoi(r)
		if err != nil {
			return nil, ErrMalformedValue
		}
	}

	return ints, nil
}

// IntsPanic returns all integer values associated with the given key or panics
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntsPanic(key string) []int {
	if value, err := this.IntsError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// IntsFatal returns all integer values associated with the given key or calls log.Fatal()
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntsFatal(key string) []int {
	if value, err := this.IntsError(key); err != nil {
		this.fatal(key, err)
		return nil
	} else {
		return value
	}
}

// IntsDefault returns all integer values associated with the given key or returns provided defaults
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntsDefault(key string, Default []int) []int {
	if value, err := this.IntsError(key); err != nil {
		return Default
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// Int returns the first integer value associated with the given key or returns 0
// if the key does not exist.
func (this *Reader) Int(key string) int {
	value, _ := this.IntError(key)
	return value
}

// IntError returns the first integer value associated with the given key with an error
// if the key does not exist or the values could not be parsed as integers (according to strconv.Atoi).
func (this *Reader) IntError(key string) (int, error) {
	raw, err := this.StringError(key)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(raw)
	if err != nil {
		return 0, ErrMalformedValue
	}

	return number, nil
}

// IntPanic returns the first integer value associated with the given key or panics
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntPanic(key string) int {
	if value, err := this.IntError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// IntFatal returns the first integer value associated with the given key or calls log.Fatal()
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntFatal(key string) int {
	if value, err := this.IntError(key); err != nil {
		this.fatal(key, err)
		return 0
	} else {
		return value
	}
}

// IntDefault returns the first integer values associated with the given key or returns the provided default
// if the key does not exist or the values could not be parsed as integers.
func (this *Reader) IntDefault(key string, Default int) int {
	if value, err := this.IntError(key); err != nil {
		return Default
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// Bool returns the boolean value associated with the given key or false
// if the key does not exist or the value could not be parsed as a bool.
func (this *Reader) Bool(key string) bool {
	value, _ := this.BoolError(key)
	return value
}

// BoolError returns the boolean value associated with the given key with an error
// if the key does not exist or the value could not be parsed as a bool (according to strconv.ParseBool).
func (this *Reader) BoolError(key string) (bool, error) {
	raw, err := this.StringError(key)
	if err != nil {
		return false, err
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, ErrMalformedValue
	}

	return value, nil
}

// BoolPanic returns the boolean value associated with the given key or panics
// if the key does not exist or the value could not be parsed as a bool.
func (this *Reader) BoolPanic(key string) bool {
	if value, err := this.BoolError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// BoolFatal returns the boolean value associated with the given key or calls log.Fatal()
// if the key does not exist or the value could not be parsed as a bool.
func (this *Reader) BoolFatal(key string) bool {
	if value, err := this.BoolError(key); err != nil {
		this.fatal(key, err)
		return false
	} else {
		return value
	}
}

// BoolDefault returns the boolean value associated with the given key or returns the provided default
// if the key does not exist or the value could not be parsed as a bool.
func (this *Reader) BoolDefault(key string, Default bool) bool {
	if value, err := this.BoolError(key); err != nil {
		return Default
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// URLs returns all URL values associated with the given key or returns the zero value
// if the key does not exist or the value could not be parsed as a URL.
func (this *Reader) URLs(key string) []url.URL {
	value, _ := this.URLsError(key)
	return value
}

// URLsError returns all URL values associated with the given key with an error
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLsError(key string) ([]url.URL, error) {
	raw, err := this.StringsError(key)
	if err != nil {
		return nil, err
	}

	urls := make([]url.URL, len(raw))
	for i, r := range raw {
		parsed, err := url.Parse(r)
		if err != nil {
			return nil, ErrMalformedValue
		}
		urls[i] = *parsed
	}

	return urls, nil
}

// URLsPanic returns all URL values associated with the given key or panics
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLsPanic(key string) []url.URL {
	if value, err := this.URLsError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// URLsFatal returns all URL values associated with the given key or calls log.Fatal()
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLsFatal(key string) []url.URL {
	if value, err := this.URLsError(key); err != nil {
		this.fatal(key, err)
		return nil
	} else {
		return value
	}
}

// URLsDefault returns all URL values associated with the given key or returns provided defaults
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLsDefault(key string, Default []url.URL) []url.URL {
	if value, err := this.URLsError(key); err != nil {
		return Default
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// URL returns the first URL associated with the given key or returns the zero value
// if the key does not exist or the value could not be parsed as a URL.
func (this *Reader) URL(key string) url.URL {
	value, _ := this.URLError(key)
	return value
}

// URLError returns the first URL associated with the given key with an error
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLError(key string) (url.URL, error) {
	raw, err := this.StringError(key)
	if err != nil {
		return url.URL{}, err
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return url.URL{}, ErrMalformedValue
	}

	return *parsed, nil
}

// URLPanic returns the first URL associated with the given key or panics
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLPanic(key string) url.URL {
	if value, err := this.URLError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

// URLFatal returns the first URL associated with the given key or calls log.Fatal()
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLFatal(key string) url.URL {
	if value, err := this.URLError(key); err != nil {
		this.fatal(key, err)
		return url.URL{}
	} else {
		return value
	}
}

// URLDefault returns the first URL associated with the given key or returns provided defaults
// if the key does not exist or the values could not be parsed as URLs.
func (this *Reader) URLDefault(key string, Default url.URL) url.URL {
	if value, err := this.URLError(key); err != nil {
		return Default
	} else {
		return value
	}
}

////////////////////////////////////////////////////////////////////////////////

// Duration returns the first Duration associated with the given key or returns the zero value
// if the key does not exist or the value could not be parsed as a Duration.
// For examples of duration strings see http://golang.org/pkg/time/#ParseDuration
func (this *Reader) Duration(key string) time.Duration {
	duration, _ := this.DurationError(key)
	return duration
}

// DurationError returns the first Duration associated with the given key with an error
// if the key does not exist or the values could not be parsed as Durations.
func (this *Reader) DurationError(key string) (time.Duration, error) {
	raw, err := this.StringError(key)
	if err != nil {
		return 0, err
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return 0, ErrMalformedValue
	}

	return parsed, nil
}

// DurationPanic returns the first Duration associated with the given key or panics
// if the key does not exist or the values could not be parsed as Durations.
func (this *Reader) DurationPanic(key string) time.Duration {
	if duration, err := this.DurationError(key); err != nil {
		panic(err)
	} else {
		return duration
	}
}

// DurationFatal returns the first Duration associated with the given key or calls log.Fatal()
// if the key does not exist or the values could not be parsed as Durations.
func (this *Reader) DurationFatal(key string) time.Duration {
	if duration, err := this.DurationError(key); err != nil {
		this.fatal(key, err)
		return 0
	} else {
		return duration
	}
}

// DurationDefault returns the first Duration associated with the given key or returns provided default
// if the key does not exist or the values could not be parsed as Durations.
func (this *Reader) DurationDefault(key string, Default time.Duration) time.Duration {
	if duration, err := this.DurationError(key); err != nil {
		return Default
	} else {
		return duration
	}
}

////////////////////////////////////////////////////////////////////////////////

// Time returns the first Time associated with the given key or returns the zero value
// if the key does not exist or the value could not be parsed as a Time using the provided format.
// For examples of format strings see http://golang.org/pkg/time/#pkg-constants
func (this *Reader) Time(key string, format string) time.Time {
	parsed, _ := this.TimeError(key, format)
	return parsed
}

// TimeError returns the first Time associated with the given key with an error
// if the key does not exist or the values could not be parsed as Times using the provided format.
func (this *Reader) TimeError(key string, format string) (time.Time, error) {
	raw, err := this.StringError(key)
	if err != nil {
		return time.Time{}, err
	}

	parsed, err := time.Parse(format, raw)
	if err != nil {
		return time.Time{}, ErrMalformedValue
	}

	return parsed, nil
}

// TimePanic returns the first Time associated with the given key or panics
// if the key does not exist or the values could not be parsed as Times using the provided format.
func (this *Reader) TimePanic(key string, format string) time.Time {
	if instant, err := this.TimeError(key, format); err != nil {
		panic(err)
	} else {
		return instant
	}
}

// TimeFatal returns the first Time associated with the given key or calls log.Fatal()
// if the key does not exist or the values could not be parsed as Times using the provided format.
func (this *Reader) TimeFatal(key string, format string) time.Time {
	if instant, err := this.TimeError(key, format); err != nil {
		this.fatal(key, err)
		return time.Time{}
	} else {
		return instant
	}
}

// TimeDefault returns the first Time associated with the given key or returns provided default
// if the key does not exist or the values could not be parsed as Times using the provided format.
func (this *Reader) TimeDefault(key string, format string, Default time.Time) time.Time {
	if instant, err := this.TimeError(key, format); err != nil {
		return Default
	} else {
		return instant
	}
}

////////////////////////////////////////////////////////////////////////////////
