package configo

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestJSONSourceFixture(t *testing.T) {
	gunit.Run(new(JSONSourceFixture), t)
}

type JSONSourceFixture struct {
	*gunit.Fixture
}

func (this *JSONSourceFixture) Setup() {}

func (this *JSONSourceFixture) TestParseMalformedJSONPanics() {
	malformed := []byte(`{"key":"value",}}}`)
	this.So(func() { FromJSONContent(malformed) }, should.Panic)
}
func (this *JSONSourceFixture) TestNonExistentSingleValue() {
	this.assertFailure(`{}`, "key", ErrKeyNotFound)
}

func (this *JSONSourceFixture) TestReadSingleStringValue() {
	this.assertSuccess(`{"key":"value"}`, "key", "value")
}
func (this *JSONSourceFixture) TestReadMultipleStringValues() {
	this.assertSuccess(`{"key": [ "a", "b", "c" ] }`, "key", "a", "b", "c")
}

func (this *JSONSourceFixture) TestReadSingleNumericValue() {
	this.assertSuccess(`{"key":1234}`, "key", "1234")
}
func (this *JSONSourceFixture) TestReadMultileNumericValues() {
	this.assertSuccess(`{"key":[1,2,3]}`, "key", "1", "2", "3")
}

func (this *JSONSourceFixture) TestReadSingleDecimalValue() {
	this.assertSuccess(`{"key":1234.5678}`, "key", "1234.5678")
}
func (this *JSONSourceFixture) TestReadMultipleDecimalValues() {
	this.assertSuccess(`{"key":[1.2, 3.4, 5.6]}`, "key", "1.2", "3.4", "5.6")
}

func (this *JSONSourceFixture) TestReadSingleBooleanValue() {
	this.assertSuccess(`{"key":true}`, "key", "true")
}
func (this *JSONSourceFixture) TestReadMultipleBooleanValues() {
	this.assertSuccess(`{"key":[true, false, true]}`, "key", "true", "false", "true")
}

func (this *JSONSourceFixture) TestReadMultipleValues() {
	this.assertSuccess(`{"key":["value", 1, 1.2, true]}`, "key", "value", "1", "1.2", "true")
}

func (this *JSONSourceFixture) assertSuccess(raw, key string, expectedValues ...string) {
	source := FromJSONContent([]byte(raw))

	values, err := source.Strings(key)

	this.So(values, should.Resemble, expectedValues)
	this.So(err, should.BeNil)
}
func (this *JSONSourceFixture) assertFailure(raw, key string, expectedError error) {
	source := FromJSONContent([]byte(raw))

	values, err := source.Strings(key)

	this.So(values, should.BeEmpty)
	this.So(err, should.Equal, expectedError)
}
