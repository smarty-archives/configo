package configo

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// JSONSource houses key-value pairs unmarshaled from JSON data.
type JSONSource struct {
	values map[string]interface{}
}

// FromJSONFile reads and unmarshals the file at the provided path into a JSONSource.
// Any resulting error results in a panic.
func FromJSONFile(filename string) Source {
	if contents, err := ioutil.ReadFile(filename); err != nil {
		panic(err)
	} else {
		return FromJSONContent(contents)
	}
}

// If the provided condition returns true, the specified filename
// is required and must be found; otherwise loading the file is optional.
func FromConditionalJSONFile(filename string, condition func() bool) Source {
	if condition() {
		return FromJSONFile(filename)
	}

	return FromOptionalJSONFile(filename)
}

// FromOptionalJSONFile is like FromJSONFile but it does not panic if the file is not found.
func FromOptionalJSONFile(filename string) Source {
	if contents, _ := ioutil.ReadFile(filename); len(contents) > 0 {
		return FromJSONContent(contents)
	}

	return NoopSource{}
}

// FromJSONContent unmarshals the provided json content into a JSONSource.
// Any resulting error results in a panic.
func FromJSONContent(raw []byte) Source {
	values := make(map[string]interface{})
	if err := json.Unmarshal(raw, &values); err != nil {
		panic(err)
	}

	return &JSONSource{values: values}
}

func (this *JSONSource) Strings(key string) ([]string, error) {
	// FUTURE: if contents of key contain a hypen, change it to an underscore?
	// FUTURE: split on / character to indicate another level
	if item, found := this.values[key]; found {
		return toStrings(item), nil
	}

	return nil, KeyNotFoundError
}

func toStrings(value interface{}) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case float64:
		return []string{strconv.FormatFloat(typed, 'f', -1, 64)}
	case bool:
		return []string{strconv.FormatBool(typed)}
	case []interface{}:
		values := []string{}
		for _, item := range typed {
			values = append(values, convertToString(item))
		}
		return values
	default:
		return nil
	}
}

func (this *JSONSource) Initialize() {}
