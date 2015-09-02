package configo

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type CommandLineSourceFixture struct {
	*gunit.Fixture

	source            *CommandLineSource
	originalFlagParse func()
}

var commandLineValue = "value from command line"
var noCommandLineValue = ""

func (this *CommandLineSourceFixture) Setup() {
	this.originalFlagParse = flagParse
}

func (this *CommandLineSourceFixture) Teardown() {
	flagParse = this.originalFlagParse
}

func (this *CommandLineSourceFixture) TestMatchingValueFound() {
	this.Print(`Simulates at the command line: ./app -flagName="value from command line"`)

	flagParse = this.fakeFlagParse
	this.source = FromCommandLineFlag("flagName", "This is a cool flag")
	this.source.Initialize()

	values, err := this.source.Strings("flagName")

	this.So(values, should.Resemble, []string{commandLineValue})
	this.So(err, should.BeNil)
}

func (this *CommandLineSourceFixture) TestFlagNotPassed__NotFound() {
	this.Print(`Simulates no command line flag passed.`)

	this.source = FromCommandLineFlag("flagName2", "This is a cool flag")
	this.source.Initialize()

	values, err := this.source.Strings("flagName2")

	var expectedValue []string // nil
	this.So(values, should.Resemble, expectedValue)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *CommandLineSourceFixture) TestFlagNotDefined__NoValuesReturned() {
	this.Print(`Simulates requesting the value of an undefined flag`)

	flagParse = this.fakeFlagParse
	this.source = FromCommandLineFlag("flagName3", "This is a cool flag")
	this.source.Initialize()

	values, err := this.source.Strings("unknown")

	var expectedValue []string // nil
	this.So(values, should.Resemble, expectedValue)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *CommandLineSourceFixture) fakeFlagParse() {
	this.source.value = &commandLineValue
	this.source.isSet = true
}
