package configo

import (
	"net/url"
	"testing"
	"time"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDefaultSourceFixture(t *testing.T) {
	gunit.Run(new(DefaultSourceFixture), t)
}

type DefaultSourceFixture struct {
	*gunit.Fixture

	source *DefaultSource
	pairs  []DefaultPair
}

func (this *DefaultSourceFixture) TestEmptyKeyReportsNoValues() {
	this.assertError(ErrKeyNotFound)
}

func (this *DefaultSourceFixture) TestAddedStringCanBeRetrieved() {
	this.addValues("Hello, World!")

	this.assertValues([]string{"Hello, World!"})
}

func (this *DefaultSourceFixture) TestAddingMultipleStringsRetrievesAll() {
	this.addValues("Hello,")
	this.addValues("World!")

	this.assertValues([]string{"Hello,", "World!"})
}

func (this *DefaultSourceFixture) TestAddingBooleanTypesRetrievesAll() {
	this.addValues(true, false)

	this.assertValues([]string{"true", "false"})
}

func (this *DefaultSourceFixture) TestAddingIntegerTypesRetrievesAll() {
	this.addValues(-1)
	this.addValues(-int64(2))
	this.addValues(-int32(3))
	this.addValues(-int16(4))
	this.addValues(-int8(5))

	this.assertValues([]string{"-1", "-2", "-3", "-4", "-5"})
}

func (this *DefaultSourceFixture) TestAddingUnsignedIntegerTypesRetrievesAll() {
	this.addValues(uint64(6))
	this.addValues(uint32(7))
	this.addValues(uint16(8))
	this.addValues(uint8(9))
	this.addValues(uint(10))

	this.assertValues([]string{"6", "7", "8", "9", "10"})
}

func (this *DefaultSourceFixture) TestAddingUnknownTypeReturnsEmptyValue() {
	this.addValues([]string{"slice"})
	this.assertValues([]string{""})
}

func (this *DefaultSourceFixture) TestAddingFloatingPointTypesRetrievesAll() {
	this.addValues(float32(1.23456), float64(1.2345678901234567))

	this.assertValues([]string{"1.23456", "1.2345678901234567"})
}

func (this *DefaultSourceFixture) TestAddingURLTypesRetrievesAll() {
	url, _ := url.Parse("https://user:pass@host:1234/path?query=value#segment")
	this.addValues(url, *url)

	this.assertValues([]string{"https://user:pass@host:1234/path?query=value#segment", "https://user:pass@host:1234/path?query=value#segment"})
}

func (this *DefaultSourceFixture) TestAddingDurationTypesRetrievesAll() {
	this.addValues(time.Second)
	this.addValues(time.Millisecond)
	this.addValues(time.Microsecond)
	this.addValues(time.Nanosecond)

	this.assertValues([]string{"1s", "1ms", "1µs", "1ns"})
}

func (this *DefaultSourceFixture) TestAddingTimeTypesRetrievesAll() {
	now := time.Now().UTC()
	this.addValues(now)

	this.assertValues([]string{now.String()})
}

func (this *DefaultSourceFixture) addValues(values ...interface{}) {
	this.pairs = append(this.pairs, Default("key", values...))
}

func (this *DefaultSourceFixture) assertValues(expected []string) {
	this.source = NewDefaultSource(this.pairs...)

	values, err := this.source.Strings("key")

	this.So(err, should.BeNil)
	this.So(values, should.Resemble, expected)

}
func (this *DefaultSourceFixture) assertError(expected error) {
	this.source = NewDefaultSource(this.pairs...)

	values, err := this.source.Strings("key")

	this.So(err, should.Equal, expected)
	this.So(values, should.BeNil)
}
