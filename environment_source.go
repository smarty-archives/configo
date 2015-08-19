package configo

import (
	"os"
	"strings"
	"unicode"
)

// EnvironmentSource reads key-value pairs from the environment.
type EnvironmentSource struct {
	prefix    string
	separator string
}

// FromEnvironment creates an envirnoment source capable of
// parsing values separated by the pipe character.
func FromEnvironment() *EnvironmentSource {
	return FromEnvironmentCustomSeparator("", "|")
}

// FromEnvironmentWithPrefix creates an envirnoment source capable of:
// - reading values with keys all beginning with the provided prefix,
// - parsing values separated by the pipe character.
func FromEnvironmentWithPrefix(prefix string) *EnvironmentSource {
	return FromEnvironmentCustomSeparator(prefix, "|")
}

// FromEnvironmentWithPrefix creates an envirnoment source capable of
// parsing values separated by the specified character.
func FromEnvironmentCustomSeparator(prefix, separator string) *EnvironmentSource {
	return &EnvironmentSource{prefix: prefix, separator: separator}
}

func (this *EnvironmentSource) Strings(key string) ([]string, error) {
	key = this.prefix + sanitizeKey(key)

	if value := os.Getenv(key); len(value) > 0 {
		return strings.Split(value, this.separator), nil
	}

	if value := os.Getenv(strings.ToUpper(key)); len(value) > 0 {
		return strings.Split(value, this.separator), nil
	}

	if value := os.Getenv(strings.ToLower(key)); len(value) > 0 {
		return strings.Split(value, this.separator), nil
	}

	return nil, KeyNotFoundError
}
func sanitizeKey(key string) string {
	sanitized := ""

	for _, character := range key {
		if unicode.IsDigit(character) {
			sanitized += string(character)
		} else if unicode.IsLetter(character) {
			sanitized += string(character)
		} else {
			sanitized += "_"
		}
	}

	return sanitized
}

func (this *EnvironmentSource) Initialize() {}
