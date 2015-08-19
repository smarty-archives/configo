package configo

import "flag"

// CommandLineSource registers a single command line flag
// and stores it's actual value, if supplied on the command line.
// It implements the Source interface so it can be used by a Reader.
type CommandLineSource struct {
	isSet bool
	name  string
	value *string
}

// NewCommandLineFlag receives the name, defaultValue, and description of a command line flag.
// The default value can be of any type handled by the internal convertString function.
func NewCommandLineFlag(name string, description string) *CommandLineSource {
	return &CommandLineSource{
		name:  name,
		value: flag.String(name, "<undefined>", description),
	}
}

// Initialize calls flag.Parse(). Do not call until all CommandLineSource instances have been created.
func (this *CommandLineSource) Initialize() {
	flagParse()

	flag.Visit(this.visitor)
}

func (this *CommandLineSource) visitor(f *flag.Flag) {
	if f.Name == this.name {
		this.isSet = true
	}
}

// Strings returns the command line flag value, or the default if no value was provided at the command line.
func (this *CommandLineSource) Strings(key string) ([]string, error) {
	if key != this.name {
		return nil, KeyNotFoundError
	}
	if !this.isSet {
		return nil, KeyNotFoundError
	}
	return []string{*this.value}, nil
}

// flagParse forwards to flag.Parse() in production but allows tests to use their own implementation.
var flagParse = func() { flag.Parse() }
