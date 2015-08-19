package configo

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ConditionalSourceFixture struct {
	*gunit.Fixture

	source *ConditionalSource
	active bool
}

func (this *ConditionalSourceFixture) Setup() {
	this.source = NewConditionalSource(func() bool { return this.active })
	this.active = true
}

func (this *ConditionalSourceFixture) TestEmptyKeyReportsNoValues() {
	this.assertError(KeyNotFoundError)
}

func (this *ConditionalSourceFixture) TestAddingMultipleStringsRetrievesAll() {
	this.addValues("Hello,")
	this.addValues("World!")

	this.assertValues([]string{"Hello,", "World!"})
}

func (this *ConditionalSourceFixture) TestFalseConditionReportsNoValues() {
	this.addValues("Hello, World!")
	this.active = false

	this.assertError(KeyNotFoundError)
}

func (this *ConditionalSourceFixture) addValues(values ...interface{}) {
	this.source.Add("key", values...)
}

func (this *ConditionalSourceFixture) assertValues(expected []string) {
	values, err := this.source.Strings("key")

	this.So(err, should.BeNil)
	this.So(values, should.Resemble, expected)

}
func (this *ConditionalSourceFixture) assertError(expected error) {
	values, err := this.source.Strings("key")

	this.So(err, should.Equal, expected)
	this.So(values, should.BeNil)
}
