// Package configo provides flexible configurtion for go applications
// from any number of 'sources', each of which implement the Source interface
// also provided by this package.
//
// Sources provided in this package:
//
//     - JSONSource (key/value pairs in JSON content)
//     - EnvironmentSource (key/value pairs from the environment)
//     - CommandLineSource (key/value pairs via command line flags)
//     - DefaultSource (key/value pairs manually configured by the application)
//     - ConditionalSource (filters key/value retreival based on a condition)
//
// These basic sources have been composed into additional sources, made
// available via the following constructor methods:
//
//     - FromDevelopmentOnlyDefaults()
//     - FromRequiredInProductionJSONFile()
//     - NewDefaultCommandLineConfigFileSource()
//     - NewCommandLineConfigFileSource(path string)
//
// Any of these sources may be provided to a Reader which is then used to
// retrieve configuration values based on keys contained in the sources.
//
// The reader can fetch values of various types:
//
//     func (*Reader) Strings(key string) []string
//     func (*Reader) String(key string) string
//     func (*Reader) Ints(key string) []int
//     func (*Reader) Int(key string) int
//     func (*Reader) Bool(key string) bool
//     func (*Reader) URLs(key string) []net.url.URL
//     func (*Reader) URL(key string) net.url.URL
//
// For each of the types returned above there are different ways to handle
// the scenario when a key is not found. I'll illustrate this with the
// applicable Int functions (but similar methods are implemented for each
// returned type):
//
//     // returns zero value if key not found or values are malformed.
//     func (*Reader) Int(key string) int
//
//     // returns the value or the specified default if the key is not found.
//     func (*Reader) IntDefault(key string, Default int) int
//
//     // returns 0 and an error if key not found or values are malformed.
//     func (*Reader) IntError(key string) (int, error)
//
//     // returns the value or panics if the key is not found or the values are malformed.
//     func (*Reader) IntPanic(key string) int
//
//     // returns the value or calls log.Fatal() if the key is not found or the values are malformed.
//     func (*Reader) IntFatal(key string) int
//
// Here's a full example:
//
//     reader := configo.NewReader(
//         NewDefaultCommandLineConfigFileSource(),
//         NewCommandLineFlag("s3-storage-address", "The address of the s3 bucket"),
//         FromOptionalJSONFile("config-prod.json"),
//     )
//     value := reader.URL("s3-storage-address")
package configo
