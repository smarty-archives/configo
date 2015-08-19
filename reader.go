package newton

import (
	"log"
	"strconv"
)

// Reader retrieves values from the provided sources, handling conversions
// to the type identified by the method being called (Strings, Ints, etc...).
type Reader struct {
	sources []Source
}

func NewReader(sources ...Source) *Reader {
	return &Reader{sources: sources}
}

// Strings returns all values associated with the given key or nil
// if the key does not exist.
func (this *Reader) Strings(key string) []string {
	value, _ := this.StringsError(key)
	return value
}

// StringsError returns all values associated with the given key with an error
// if the key does not exist.
func (this *Reader) StringsError(key string) ([]string, error) {
	for _, source := range this.sources {
		if value, err := source.Strings(key); err == nil {
			return value, nil
		}
	}
	return nil, KeyNotFoundError
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
		fatal(err)
		return nil
	} else {
		return value
	}
}

//////////////////////////////////////////////////////////////

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
		fatal(err)
		return ""
	} else {
		return value
	}
}

//////////////////////////////////////////////////////////////

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
			return nil, MalformedValueError
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
		fatal(err)
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

//////////////////////////////////////////////////////////////

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
		return 0, MalformedValueError
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
		fatal(err)
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

//////////////////////////////////////////////////////////////

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
		return false, MalformedValueError
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
		fatal(err)
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

//////////////////////////////////////////////////////////////

var fatal = func(err error) { log.Fatal(err) }
