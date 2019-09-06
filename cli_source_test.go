package configo

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCLISourceFixture(t *testing.T) {
	gunit.Run(new(CLISourceFixture), t)
}

type CLISourceFixture struct {
	*gunit.Fixture
	source *CLISource
	output *bytes.Buffer
}

var commandLineValue = "value from command line"

func (this *CLISourceFixture) TestMatchingValueFound() {
	this.Println(`Simulates at the command line: ./app -flagName="value from command line"`)

	this.source = FromCLI(Flag("flagName", "This is a cool flag"))
	this.source.source = []string{"./app", fmt.Sprintf("-flagName=%s", commandLineValue)}
	this.source.Initialize()

	values, err := this.source.Strings("flagName")

	this.So(values, should.Resemble, []string{commandLineValue})
	this.So(err, should.BeNil)
}

func (this *CLISourceFixture) TestFlagNotPassed__NotFound() {
	this.Println(`Simulates no command line flag passed.`)

	this.source = FromCLI(Flag("flagName", "This is a cool flag"))
	this.source.source = []string{"./app"}
	this.source.Initialize()

	values, err := this.source.Strings("flagName")

	var expectedValue []string // nil
	this.So(values, should.Resemble, expectedValue)
	this.So(err, should.Equal, ErrKeyNotFound)
}

func (this *CLISourceFixture) TestFlagNotDefined__NoValuesReturned() {
	this.Println(`Simulates requesting the value of an undefined flag`)

	this.source = FromCLI(Flag("flagName", "This is a cool flag"))
	this.source.source = []string{"./app"}
	this.source.Initialize()

	values, err := this.source.Strings("unknown")

	var expectedValue []string // nil
	this.So(values, should.Resemble, expectedValue)
	this.So(err, should.Equal, ErrKeyNotFound)
}

func (this *CLISourceFixture) TestBooleanFlagDefined() {
	this.Println(`Simulates requesting the value of a boolean flag`)
	this.source = FromCLI(
		BoolFlag("flagname1", "This is cool"),
		BoolFlag("flagname2", "This is cooler"),
		BoolFlag("flagname3", "This is coolest"),
		BoolFlag("flagname4", "This is stellar"),
	)
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
	this.So(err, should.Equal, ErrKeyNotFound)

	values, err = this.source.Strings("flagname5")
	this.So(values, should.BeNil)
	this.So(err, should.Equal, ErrKeyNotFound)
}

func (this *CLISourceFixture) TestUsageMessage() {
	buffer := new(bytes.Buffer)
	this.source = FromCLI(
		Flag("something", "yeah"),
		Usage("This is some helpful text"),
		ContinueOnError(),
		SetOutput(buffer),
	)
	this.source.source = []string{"./app", "-help", "-something"}

	this.source.Initialize()

	this.So(buffer.String(), should.ContainSubstring, "-something")
	this.So(buffer.String(), should.ContainSubstring, "yeah")
	this.So(buffer.String(), should.ContainSubstring, "This is some helpful text")
}
