package configo

import (
	"flag"
	"os"
)

// CommandLineSource allows for registration of command line flags
// and stores their actual values, if supplied on the command line.
// It implements the Source interface so it can be used by a Reader.
type CommandLineSource struct {
	source   []string
	flags    *flag.FlagSet
	registry map[string]*string
	values   map[string]string
}

const flagSetName = "configo"

// FromCommandLineFlags creates a new CommandLineSource for use in a Reader.
// It uses a *flag.FlagSet internally to register and parse the flags.
// Be default the flag.ErrorHandling mode is set to flag.ExitOnError
func FromCommandLineFlags() *CommandLineSource {
	return &CommandLineSource{
		source:   os.Args,
		flags:    flag.NewFlagSet(flagSetName, flag.ExitOnError),
		registry: make(map[string]*string),
		values:   make(map[string]string),
	}
}

// ContinueOnError sets the flag.ErrorHandling mode of the internal *flag.FlagSet
// to flag.ContinueOnError. Must be called before Initialize is called.
func (this *CommandLineSource) ContinueOnError() *CommandLineSource {
	this.flags.Init(flagSetName, flag.ContinueOnError)
	return this
}

// ContinueOnError sets the flag.ErrorHandling mode of the internal *flag.FlagSet
// to flag.PanicOnError. Must be called before Initialize is called.
func (this *CommandLineSource) PanicOnError() *CommandLineSource {
	this.flags.Init(flagSetName, flag.PanicOnError)
	return this
}

// Register adds flags and corresponding usage descriptions to the CommandLineSource.
func (this *CommandLineSource) Register(name, description string) *CommandLineSource {
	this.registry[name] = this.flags.String(name, "", description)
	return this
}

// Parses the internal *flag.FlagSet. Call only after making all Register calls.
func (this *CommandLineSource) Initialize() {
	this.flags.Parse(this.source[1:])
	this.flags.Visit(this.visitor)
}

func (this *CommandLineSource) visitor(flag *flag.Flag) {
	this.values[flag.Name] = *this.registry[flag.Name]
}

// Strings returns the matching command line flag value, or KeyNotFound.
func (this *CommandLineSource) Strings(key string) ([]string, error) {
	value, found := this.values[key]
	if !found {
		return nil, KeyNotFoundError
	}
	return []string{value}, nil
}
