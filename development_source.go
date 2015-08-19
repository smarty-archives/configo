package newton

import (
	"os"
	"runtime"
)

var developmentEnvironment = func() bool {
	hostname, _ := os.Hostname()
	return runtime.GOOS == "darwin" || hostname == "vagrant"
}

var productionEnvironment = func() bool {
	return !developmentEnvironment()
}

// A conditional source which determine if we are running
// in a development environment
func NewDevelopmentSource() *ConditionalSource {
	return NewConditionalSource(developmentEnvironment)
}

// If we are running in a production (non-development) environment
// the specified JSON filename is required; otherwise it's optional.
func FromProductionJSONFile(filename string) *JSONSource {
	return FromConditionalJSONFile(filename, productionEnvironment)
}
