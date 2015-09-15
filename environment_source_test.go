package configo

import (
	"os"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type EnvironmentSourceFixture struct {
	*gunit.Fixture

	source *EnvironmentSource
}

func (this *EnvironmentSourceFixture) Setup() {
	this.source = FromEnvironmentCustomSeparator("configo_", ",")
}

func (this *EnvironmentSourceFixture) TestNonExistentValue() {
	values, err := this.source.Strings("notfound")

	this.So(values, should.BeEmpty)
	this.So(err, should.Equal, KeyNotFoundError)
}
func (this *EnvironmentSourceFixture) TestKnownValue() {
	setEnvironment("configo_Found_Single", "hello")
	values, err := this.source.Strings("Found_Single")

	this.So(values, should.Resemble, []string{"hello"})
	this.So(err, should.BeNil)
}
func (this *EnvironmentSourceFixture) TestKnownValueArray() {
	setEnvironment("configo_Found_Multiple", "a,b,c")
	values, err := this.source.Strings("Found_Multiple")

	this.So(values, should.Resemble, []string{"a", "b", "c"})
	this.So(err, should.BeNil)
}
func (this *EnvironmentSourceFixture) TestChecksUpperCaseKey() {
	setEnvironment("CONFIGO_UPPERCASE", "value")
	values, err := this.source.Strings("uppercase")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}

func (this *EnvironmentSourceFixture) TestChecksLowerCaseKey() {
	setEnvironment("configo_lowercase", "value")
	values, err := this.source.Strings("LOWERCASE")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}

func (this *EnvironmentSourceFixture) TestInvalidCharactersReplacedWithUnderscore() {
	setEnvironment("configo_0_0_0_0_0_0_0_0_0_0_0_0_0_0_0", "value")
	values, err := this.source.Strings("0-0.0~0!0@0#0%0^0&0*0(0)0/0\\0")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}

func (this *EnvironmentSourceFixture) TestEnvPrefixIsAlwaysStrippedBeforeTheEnvironmentLookup() {
	setEnvironment("configo_my_awesome_variable", "my_awesome_value")
	values, err := this.source.Strings("env:my_awesome_variable")

	this.So(values, should.Resemble, []string{"my_awesome_value"})
	this.So(err, should.BeNil)
}

func setEnvironment(key, value string) {
	os.Setenv(key, value)
}
