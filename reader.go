package newton

import (
	"log"
	"strconv"
)

type Reader struct {
	sources []Source
}

func NewReader(sources ...Source) *Reader {
	return &Reader{sources: sources}
}

func (this *Reader) Strings(key string) []string {
	value, _ := this.StringsError(key)
	return value
}

func (this *Reader) StringsError(key string) ([]string, error) {
	for _, source := range this.sources {
		if value, err := source.Strings(key); err == nil {
			return value, nil
		}
	}
	return nil, KeyNotFoundError
}

func (this *Reader) StringsPanic(key string) []string {
	if value, err := this.StringsError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func (this *Reader) StringsDefault(key string, Default []string) []string {
	if value, err := this.StringsError(key); err != nil {
		return Default
	} else {
		return value
	}
}

func (this *Reader) StringsFatal(key string) []string {
	if value, err := this.StringsError(key); err != nil {
		fatal(err)
		return nil
	} else {
		return value
	}
}

//////////////////////////////////////////////////////////////

func (this *Reader) String(key string) string {
	if value, err := this.StringError(key); err != nil {
		return ""
	} else {
		return value
	}
}

func (this *Reader) StringError(key string) (string, error) {
	if value, err := this.StringsError(key); err != nil {
		return "", err
	} else if len(value) == 0 {
		return "", nil
	} else {
		return value[0], nil
	}
}

func (this *Reader) StringPanic(key string) string {
	if value, err := this.StringError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func (this *Reader) StringDefault(key string, Default string) string {
	if value, err := this.StringError(key); err != nil {
		return Default
	} else {
		return value
	}
}

func (this *Reader) StringFatal(key string) string {
	if value, err := this.StringError(key); err != nil {
		fatal(err)
		return ""
	} else {
		return value
	}
}

//////////////////////////////////////////////////////////////

func (this *Reader) Ints(key string) []int {
	value, _ := this.IntsError(key)
	return value
}

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

func (this *Reader) IntsPanic(key string) []int {
	if value, err := this.IntsError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func (this *Reader) IntsFatal(key string) []int {
	if value, err := this.IntsError(key); err != nil {
		fatal(err)
		return nil
	} else {
		return value
	}
}

func (this *Reader) IntsDefault(key string, Default []int) []int {
	if value, err := this.IntsError(key); err != nil {
		return Default
	} else {
		return value
	}
}

//////////////////////////////////////////////////////////////

func (this *Reader) Int(key string) int {
	value, _ := this.IntError(key)
	return value
}

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

func (this *Reader) IntPanic(key string) int {
	if value, err := this.IntError(key); err != nil {
		panic(err)
	} else {
		return value
	}
}

func (this *Reader) IntFatal(key string) int {
	if value, err := this.IntError(key); err != nil {
		fatal(err)
		return 0
	} else {
		return value
	}
}

func (this *Reader) IntDefault(key string, Default int) int {
	if value, err := this.IntError(key); err != nil {
		return Default
	} else {
		return value
	}
}

//////////////////////////////////////////////////////////////

var fatal = func(err error) { log.Fatal(err) }
