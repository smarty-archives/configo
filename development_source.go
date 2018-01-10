package configo

import (
	"os"
	"os/user"
	"runtime"
)

var developmentEnvironment = func() bool {
	hostname, _ := os.Hostname()
	vagrant, _ := user.Lookup("vagrant")
	return runtime.GOOS == "darwin" ||
		hostname == "vagrant" ||
		hostname == "ubuntu1404" ||
		vagrant != nil
}

var productionEnvironment = func() bool {
	return !developmentEnvironment()
}

// A conditional source which determines if we are running
// in a development environment
func FromDevelopmentOnlyDefaults(pairs ...DefaultPair) *ConditionalSource {
	return NewConditionalSource(developmentEnvironment, pairs...)
}

// If we are running in a production (non-development) environment
// the specified JSON filename is required; otherwise it's optional.
func FromRequiredInProductionJSONFile(filename string) Source {
	return FromConditionalJSONFile(filename, productionEnvironment)
}
