package configo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCommandLineSourceFixture(t *testing.T) {
	gunit.Run(new(CommandLineSourceFixture), t)
}

type CommandLineSourceFixture struct {
	*gunit.Fixture
	source *CommandLineSource
	output *bytes.Buffer
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

func (this *CommandLineSourceFixture) TestBooleanFlagDefined() {
	this.Println(`Simulates requesting the value of a boolean flag`)
	this.source = FromCommandLineFlags().
		RegisterBool("flagname1", "This is cool").
		RegisterBool("flagname2", "This is cooler").
		RegisterBool("flagname3", "This is coolest").
		RegisterBool("flagname4", "This is stellar")
	this.source.source = []string{"./app", "-flagname1", "-flagname2=true", "-flagname3=false"}
	this.source.Initialize()

	values, err := this.source.Strings("flagname1")
	this.So(values, should.Resemble, []string{"true"})
	this.So(err, should.BeNil)

	values, err = this.source.Strings("flagname2")
	this.So(values, should.Resemble, []string{"true"})
	this.So(err, should.BeNil)

	values, err = this.source.Strings("flagname3")
	this.So(values, should.Resemble, []string{"false"})
	this.So(err, should.BeNil)

	values, err = this.source.Strings("flagname4")
	this.So(values, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)

	values, err = this.source.Strings("flagname5")
	this.So(values, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *CommandLineSourceFixture) TestUsageMessage() {
	this.source = FromCommandLineFlags().
		Register("something", "yeah").
		Usage("This is some helpful text").
		ContinueOnError()
	this.source.source = []string{"./app", "-help", "-something"}
	buffer := new(bytes.Buffer)
	this.source.SetOutput(buffer)

	this.source.Initialize()

	this.So(buffer.String(), should.ContainSubstring, "-something")
	this.So(buffer.String(), should.ContainSubstring, "yeah")
	this.So(buffer.String(), should.ContainSubstring, "This is some helpful text")
}
