package configo

import (
	"fmt"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type CommandLineSourceFixture struct {
	*gunit.Fixture
	source *CommandLineSource
}

var commandLineValue = "value from command line"
var noCommandLineValue = ""

func (this *CommandLineSourceFixture) TestMatchingValueFound() {
	this.Println(`Simulates at the command line: ./app -flagName="value from command line"`)

	this.source = FromCommandLineFlags().Register("flagName", "This is a cool flag")
	this.source.source = []string{"./app", fmt.Sprintf("-flagName=%s", commandLineValue)}
	this.source.Initialize()

	values, err := this.source.Strings("flagName")

	this.So(values, should.Resemble, []string{commandLineValue})
	this.So(err, should.BeNil)
}

func (this *CommandLineSourceFixture) TestFlagNotPassed__NotFound() {
	this.Println(`Simulates no command line flag passed.`)

	this.source = FromCommandLineFlags().Register("flagName", "This is a cool flag")
	this.source.source = []string{"./app"}
	this.source.Initialize()

	values, err := this.source.Strings("flagName")

	var expectedValue []string // nil
	this.So(values, should.Resemble, expectedValue)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *CommandLineSourceFixture) TestFlagNotDefined__NoValuesReturned() {
	this.Println(`Simulates requesting the value of an undefined flag`)

	this.source = FromCommandLineFlags().Register("flagName", "This is a cool flag")
	this.source.source = []string{"./app"}
	this.source.Initialize()

	values, err := this.source.Strings("unknown")

	var expectedValue []string // nil
	this.So(values, should.Resemble, expectedValue)
	this.So(err, should.Equal, KeyNotFoundError)
}
