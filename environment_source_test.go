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
	this.source = FromEnvironmentCustomSeparator("newton_", ",")
}

func (this *EnvironmentSourceFixture) TestNonExistentValue() {
	values, err := this.source.Strings("notfound")

	this.So(values, should.BeEmpty)
	this.So(err, should.Equal, KeyNotFoundError)
}
func (this *EnvironmentSourceFixture) TestKnownValue() {
	setEnvironment("newton_Found", "hello")
	values, err := this.source.Strings("Found")

	this.So(values, should.Resemble, []string{"hello"})
	this.So(err, should.BeNil)
}
func (this *EnvironmentSourceFixture) TestKnownValueArray() {
	setEnvironment("newton_Found", "a,b,c")
	values, err := this.source.Strings("Found")

	this.So(values, should.Resemble, []string{"a", "b", "c"})
	this.So(err, should.BeNil)
}
func (this *EnvironmentSourceFixture) TestChecksUpperCaseKey() {
	setEnvironment("NEWTON_UPPERCASE", "value")
	values, err := this.source.Strings("uppercase")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}

func (this *EnvironmentSourceFixture) TestChecksLowerCaseKey() {
	setEnvironment("newton_lowercase", "value")
	values, err := this.source.Strings("LOWERCASE")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}

func (this *EnvironmentSourceFixture) TestInvalidCharactersReplacedWithUnderscore() {
	setEnvironment("newton_0_0_0_0_0_0_0_0_0_0_0_0_0_0_0", "value")
	values, err := this.source.Strings("0-0.0~0!0@0#0%0^0&0*0(0)0/0\\0")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}

func setEnvironment(key, value string) {
	os.Setenv(key, value)
}
