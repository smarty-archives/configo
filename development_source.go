package configo

import (
	"os"
	"runtime"
)

var developmentEnvironment = func() bool {
	hostname, _ := os.Hostname()
	return runtime.GOOS == "darwin" || hostname == "vagrant" || hostname == "ubuntu1404"
}

var productionEnvironment = func() bool {
	return !developmentEnvironment()
}

// A conditional source which determines if we are running
// in a development environment
func FromDevelopmentOnlyDefaults() *ConditionalSource {
	return NewConditionalSource(developmentEnvironment)
}

// If we are running in a production (non-development) environment
// the specified JSON filename is required; otherwise it's optional.
func FromRequiredInProductionJSONFile(filename string) Source {
	return FromConditionalJSONFile(filename, productionEnvironment)
}
