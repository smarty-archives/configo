package configo

// CLIConfigFileSource registers a command line flag for specifying an optional json config file,
// beyond any other config file definitions that follow this one. It is intended to be used to provide an
// override to the regularly used config file(s), like when you might be debugging in production (admit it,
// you've been there too).
type CLIConfigFileSource struct {
	flagName    string
	commandLine *CLISource
	json        Source
}

// FromDefaultCLIConfigFileSource registers a command line flag called "config" for specifying
// an alternate JSON config file.
func FromDefaultCLIConfigFileSource() *CLIConfigFileSource {
	return FromCLIConfigFileSource("config")
}

// FromCLIConfigFileSource registers a command line flag with the given flagName for specifying
// an alternate JSON config file.
func FromCLIConfigFileSource(flagName string) *CLIConfigFileSource {
	return &CLIConfigFileSource{
		flagName:    flagName,
		commandLine: FromCLI(Flag(flagName, "The default configuration file path.")),
	}
}

// Initialize parses the command line flag and reads the altnerate JSON source.
func (this *CLIConfigFileSource) Initialize() {
	this.commandLine.Initialize()

	path, err := this.commandLine.Strings(this.flagName)
	if err == nil && len(path) > 0 {
		this.json = FromOptionalJSONFile(path[0])
	}
}

// Strings reads the key from the JSON source if it was successfully loaded during Initialize.
func (this *CLIConfigFileSource) Strings(key string) ([]string, error) {
	if this.json == nil {
		return nil, KeyNotFoundError
	}
	return this.json.Strings(key)
}
