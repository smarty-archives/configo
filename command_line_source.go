package configo

import (
	"flag"
	"os"
	"strconv"
)

// CommandLineSource allows for registration of command line flags
// and stores their actual values, if supplied on the command line.
// It implements the Source interface so it can be used by a Reader.
type CommandLineSource struct {
	source       []string
	flags        *flag.FlagSet
	registry     map[string]*string
	boolRegistry map[string]*bool
	values       map[string]string
}

const flagSetName = "configo"

// FromCommandLineFlags creates a new CommandLineSource for use in a Reader.
// It uses a *flag.FlagSet internally to register and parse the flags.
// Be default the flag.ErrorHandling mode is set to flag.ExitOnError
func FromCommandLineFlags() *CommandLineSource {
	return &CommandLineSource{
		source:       os.Args,
		flags:        flag.NewFlagSet(flagSetName, flag.ExitOnError),
		registry:     make(map[string]*string),
		boolRegistry: make(map[string]*bool),
		values:       make(map[string]string),
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

// RegisterBool adds boolean flags and corresponding usage descriptions to the CommandLineSource.
// The advantage of this method over Register for boolean values is that the user can merely
// supply the flag without a value to set the boolean flag to true. This doesn't work with Register.
func (this *CommandLineSource) RegisterBool(name, description string) *CommandLineSource {
	this.boolRegistry[name] = this.flags.Bool(name, false, description)
	return this
}

// Parses the internal *flag.FlagSet. Call only after making all Register calls.
func (this *CommandLineSource) Initialize() {
	this.flags.Parse(this.source[1:])
	this.flags.Visit(this.visitor)
}

func (this *CommandLineSource) visitor(flag *flag.Flag) {
	if b, found := this.boolRegistry[flag.Name]; found {
		this.values[flag.Name] = strconv.FormatBool(*b)
	} else {
		this.values[flag.Name] = *this.registry[flag.Name]
	}
}

// Strings returns the matching command line flag value, or KeyNotFound.
func (this *CommandLineSource) Strings(key string) ([]string, error) {
	value, found := this.values[key]
	if !found {
		return nil, KeyNotFoundError
	}
	return []string{value}, nil
}
