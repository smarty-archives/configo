package newton

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type JSONSourceFixture struct {
	*gunit.Fixture
}

func (this *JSONSourceFixture) Setup() {}

func (this *JSONSourceFixture) TestParseMalformedJSONPanics() {
	malformed := []byte(`{"key":"value",}}}`)
	this.So(func() { FromJSONContent(malformed) }, should.Panic)
}

func (this *JSONSourceFixture) TestNonExistentSingleValue() {
	source := FromJSONContent([]byte(`{}`))

	values, err := source.Values("key")

	this.So(values, should.BeEmpty)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *JSONSourceFixture) TestReadSingleStringValue() {
	source := FromJSONContent([]byte(`{"key":"value"}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"value"})
	this.So(err, should.BeNil)
}
func (this *JSONSourceFixture) TestReadMultipleStringValues() {
	source := FromJSONContent([]byte(`{"key": [ "a", "b", "c" ] }`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"a", "b", "c"})
	this.So(err, should.BeNil)
}

func (this *JSONSourceFixture) TestReadSingleNumericValue() {
	source := FromJSONContent([]byte(`{"key":1234}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"1234"})
	this.So(err, should.BeNil)
}
func (this *JSONSourceFixture) TestReadMultipleNumericValues() {
	source := FromJSONContent([]byte(`{"key":[1,2,3]}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"1", "2", "3"})
	this.So(err, should.BeNil)
}

func (this *JSONSourceFixture) TestReadSingleDecimalValue() {
	source := FromJSONContent([]byte(`{"key":1234.5678}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"1234.5678"})
	this.So(err, should.BeNil)
}
func (this *JSONSourceFixture) TestReadMultipleDecimalValues() {
	source := FromJSONContent([]byte(`{"key":[1.2, 3.4, 5.6]}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"1.2", "3.4", "5.6"})
	this.So(err, should.BeNil)
}

func (this *JSONSourceFixture) TestReadSingleBooleanValue() {
	source := FromJSONContent([]byte(`{"key":true}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"true"})
	this.So(err, should.BeNil)
}
func (this *JSONSourceFixture) TestReadMultipleBooleanValues() {
	source := FromJSONContent([]byte(`{"key":[true, false, true]}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"true", "false", "true"})
	this.So(err, should.BeNil)
}

func (this *JSONSourceFixture) TestReadMultipleValues() {
	source := FromJSONContent([]byte(`{"key":["value", 1, 1.2, true]}`))

	values, err := source.Values("key")

	this.So(values, should.Resemble, []string{"value", "1", "1.2", "true"})
	this.So(err, should.BeNil)
}
