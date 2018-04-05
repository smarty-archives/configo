package configo

import (
	"os"
	"os/user"
	"runtime"
)

// A conditional source which determines if we are running in a development environment.
func FromDevelopmentOnlyDefaults(pairs ...DefaultPair) *ConditionalSource {
	return NewConditionalSource(developmentEnvironment, pairs...)
}
func developmentEnvironment() bool {
	return runtime.GOOS == "darwin" || hostname() == "vagrant" || hasVagrantUser()
}
func hostname() string {
	name, _ := os.Hostname()
	return name
}
func hasVagrantUser() bool {
	vagrant, _ := user.Lookup("vagrant")
	return vagrant != nil
}

// If we are running in a production (non-development) environment
// the specified JSON filename is required; otherwise it's optional.
func FromRequiredInProductionJSONFile(filename string) Source {
	return FromConditionalJSONFile(filename, productionEnvironment)
}

func productionEnvironment() bool {
	return !developmentEnvironment()
}
