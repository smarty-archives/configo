package configo

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// CLISource allows for registration of command line flags
// and stores their actual values, if supplied on the command line.
// It implements the Source interface so it can be used by a Reader.
type CLISource struct {
	source       []string
	flags        *flag.FlagSet
	registry     map[string]*string
	boolRegistry map[string]*bool
	values       map[string]string
	output       io.Writer
	usageMessage string
}

const flagSetName = "configo"

// FromCLI creates a new CLISource for use in a Reader.
// It uses a *flag.FlagSet internally to register and parse the flags.
// Be default the flag.ErrorHandling mode is set to flag.ExitOnError
func FromCLI(options ...CLI) *CLISource {
	source := &CLISource{
		source:       os.Args,
		flags:        flag.NewFlagSet(flagSetName, flag.ExitOnError),
		registry:     make(map[string]*string),
		boolRegistry: make(map[string]*bool),
		values:       make(map[string]string),
	}
	for _, option := range options {
		option(source)
	}
	return source
}

type CLI func(*CLISource)

// ContinueOnError sets the flag.ErrorHandling mode of the internal *flag.FlagSet
// to flag.ContinueOnError. Must be called before Initialize is called.
func ContinueOnError() CLI {
	return func(this *CLISource) { this.flags.Init(flagSetName, flag.ContinueOnError) }
}

// PanicOnError sets the flag.ErrorHandling mode of the internal *flag.FlagSet
// to flag.PanicOnError. Must be called before Initialize is called.
func PanicOnError() CLI {
	return func(this *CLISource) { this.flags.Init(flagSetName, flag.PanicOnError) }
}

// Flag registers a flag and corresponding usage description with the CLISource.
func Flag(name, description string) CLI {
	return func(this *CLISource) { this.registry[name] = this.flags.String(name, "", description) }
}

// BoolFlag registers a boolean flag and corresponding usage description with the CLISource.
// The advantage of this method over Flag for boolean values is that the user can merely
// supply the flag without a value to set the boolean flag to true. This doesn't work with Flag.
func BoolFlag(name, description string) CLI {
	return func(this *CLISource) { this.boolRegistry[name] = this.flags.Bool(name, false, description) }
}

// Usage appends a custom message to the end of what is normally printed
// by flag.PrintDefaults().
func Usage(message string) CLI {
	return func(this *CLISource) { this.usageMessage = message }
}

// SetOutput allows printing to an io.Writer other than os.Stderr, the default.
func SetOutput(writer io.Writer) CLI {
	return func(this *CLISource) { this.flags.SetOutput(writer); this.output = writer }
}

// Initialize parses the internal *flag.FlagSet. Call only after making all Flag calls.
func (this *CLISource) Initialize() {
	this.flags.Usage = this.usage
	this.flags.Parse(this.source[1:])
	this.flags.Visit(this.visitor)
}

func (this *CLISource) usage() {
	fmt.Fprintf(this.out(), "Usage of %s:\n", os.Args[0])
	this.flags.PrintDefaults()
	fmt.Fprintln(this.out(), this.usageMessage)
}

func (this *CLISource) out() io.Writer {
	if this.output == nil {
		return os.Stderr
	}
	return this.output
}

func (this *CLISource) visitor(flag *flag.Flag) {
	if b, found := this.boolRegistry[flag.Name]; found {
		this.values[flag.Name] = strconv.FormatBool(*b)
	} else {
		this.values[flag.Name] = *this.registry[flag.Name]
	}
}

// Strings returns the matching command line flag value, or KeyNotFound.
func (this *CLISource) Strings(key string) ([]string, error) {
	value, found := this.values[key]
	if !found {
		return nil, KeyNotFoundError
	}
	return []string{value}, nil
}
