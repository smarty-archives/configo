package configo

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestConditionalSourceFixture(t *testing.T) {
	gunit.Run(new(ConditionalSourceFixture), t)
}

type ConditionalSourceFixture struct {
	*gunit.Fixture

	source *ConditionalSource
	active bool
	pairs  []DefaultPair
}

func (this *ConditionalSourceFixture) Setup() {
	this.active = true
}

func (this *ConditionalSourceFixture) TestEmptyKeyReportsNoValues() {
	this.assertError(ErrKeyNotFound)
}

func (this *ConditionalSourceFixture) TestAddingMultipleStringsRetrievesAll() {
	this.addValues("Hello,")
	this.addValues("World!")

	this.assertValues([]string{"Hello,", "World!"})
}

func (this *ConditionalSourceFixture) TestFalseConditionReportsNoValues() {
	this.addValues("Hello, World!")
	this.active = false

	this.assertError(ErrKeyNotFound)
}

func (this *ConditionalSourceFixture) addValues(values ...interface{}) {
	this.pairs = append(this.pairs, Default("key", values...))
}

func (this *ConditionalSourceFixture) assertValues(expected []string) {
	this.source = NewConditionalSource(func() bool { return this.active }, this.pairs...)

	values, err := this.source.Strings("key")

	this.So(err, should.BeNil)
	this.So(values, should.Resemble, expected)

}
func (this *ConditionalSourceFixture) assertError(expected error) {
	this.source = NewConditionalSource(func() bool { return this.active }, this.pairs...)

	values, err := this.source.Strings("key")

	this.So(err, should.Equal, expected)
	this.So(values, should.BeNil)
}
