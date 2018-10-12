package configo

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
)

// JSONSource houses key-value pairs unmarshaled from JSON data.
type JSONSource struct {
	values map[string]interface{}
}

// FromConfigurableJSONFile allows the user to configure the config file path
// via the -config command line flag.
func FromConfigurableJSONFile() *JSONSource {
	flags := flag.NewFlagSet("config-file", flag.ContinueOnError)
	filename := flags.String("config", "config.json", "The path to the JSON config file.")
	flags.Parse(os.Args[1:]) // don't include the command name (argument #0).
	return FromJSONFile(*filename)
}

// FromJSONFile reads and unmarshals the file at the provided path into a JSONSource.
// Any resulting error results in a panic.
func FromJSONFile(filename string) *JSONSource {
	if contents, err := ioutil.ReadFile(filename); err != nil {
		panic(err)
	} else {
		return FromJSONContent(contents)
	}
}

// If the provided condition returns true, the specified filename
// is required and must be found; otherwise loading the file is optional.
func FromConditionalJSONFile(filename string, condition func() bool) *JSONSource {
	if condition() {
		return FromJSONFile(filename)
	}

	return FromOptionalJSONFile(filename)
}

// FromOptionalJSONFile is like FromJSONFile but it does not panic if the file is not found.
func FromOptionalJSONFile(filename string) *JSONSource {
	if contents, _ := ioutil.ReadFile(filename); len(contents) > 0 {
		return FromJSONContent(contents)
	}

	return nil
}

// FromJSONContent unmarshals the provided json content into a JSONSource.
// Any resulting error results in a panic.
func FromJSONContent(raw []byte) *JSONSource {
	values := make(map[string]interface{})
	if err := json.Unmarshal(raw, &values); err != nil {
		panic("json error: " + err.Error())
	}

	return FromJSONObject(values)
}

func FromJSONObject(values map[string]interface{}) *JSONSource {
	return &JSONSource{values: values}
}

func (this *JSONSource) Strings(key string) ([]string, error) {
	// FUTURE: if contents of key contain a hyphen, change it to an underscore?
	// FUTURE: split on / character to indicate another level
	if item, found := this.values[key]; found {
		return toStrings(item), nil
	}

	return nil, KeyNotFoundError
}

func toStrings(value interface{}) (values []string) {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case float64:
		return []string{strconv.FormatFloat(typed, 'f', -1, 64)}
	case bool:
		return []string{strconv.FormatBool(typed)}
	case []interface{}:
		for _, item := range typed {
			values = append(values, convertToString(item))
		}
		return values
	default:
		return nil
	}
}

func (this *JSONSource) Initialize() {}
