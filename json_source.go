package newton

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type JSONSource struct {
	values map[string]interface{}
}

func FromJSONFile(filename string) *JSONSource {
	if contents, err := ioutil.ReadFile(filename); err != nil {
		panic(err)
	} else {
		return FromJSONContent(contents)
	}
}
func FromOptionalJSONFile(filename string) *JSONSource {
	if contents, _ := ioutil.ReadFile(filename); len(contents) > 0 {
		return FromJSONContent(contents)
	}

	return nil
}

func FromJSONContent(raw []byte) *JSONSource {
	values := make(map[string]interface{})
	if err := json.Unmarshal(raw, &values); err != nil {
		panic(err)
	}

	return &JSONSource{values: values}
}

func (this *JSONSource) Name() string {
	return "json-file"
}

func (this *JSONSource) Values(key string) ([]string, error) {
	// TODO: if contents of key contain a hypen, change it to an underscore?
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
			values = append(values, toString(item))
		}
		return values
	default:
		return nil
	}
}

func toString(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return ""
	}
}
