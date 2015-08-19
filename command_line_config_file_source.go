package configo

// CommandLineConfigFileSource registers a command line flag for specifying an optional json config file,
// beyond any other config file definitions that follow this one. It is intened to be used to provide an
// orveride to the regularly used config file(s), like when you might be debugging in production (admit it,
// you've been there too).
type CommandLineConfigFileSource struct {
	flagName    string
	commandLine *CommandLineSource
	json        Source
}

// NewDefaultCommandLineConfigFileSource registers a command line flag called "config" for specifying
// an alternate JSON config file.
func NewDefaultCommandLineConfigFileSource() *CommandLineConfigFileSource {
	return NewCommandLineConfigFileSource("config")
}

// NewDefaultCommandLineConfigFileSource registers a command line flag with the given flagName for specifying
// an alternate JSON config file.
func NewCommandLineConfigFileSource(flagName string) *CommandLineConfigFileSource {
	return &CommandLineConfigFileSource{
		flagName:    flagName,
		commandLine: NewCommandLineFlag(flagName, "The default configuration file path."),
	}
}

// Initialize parses the command line flag and reads the altnerate JSON source.
func (this *CommandLineConfigFileSource) Initialize() {
	this.commandLine.Initialize()

	path, err := this.commandLine.Strings(this.flagName)
	if err == nil && len(path) > 0 {
		this.json = FromOptionalJSONFile(path[0])
	}
}

// Strings reads the key from the JSON source if it was successfully loaded during Initialize.
func (this *CommandLineConfigFileSource) Strings(key string) ([]string, error) {
	if this.json == nil {
		return nil, KeyNotFoundError
	}
	return this.json.Strings(key)
}
