package newton

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type ReaderTestFixture struct {
	*gunit.Fixture

	sources []Source
	reader  *Reader
}

func (this *ReaderTestFixture) Setup() {
	this.sources = []Source{
		&FakeSource{},
		&FakeSource{key: "string", value: []string{"asdf"}},
		&FakeSource{key: "string", value: []string{"qwer"}},
		&FakeSource{key: "int", value: []string{"42"}},
		&FakeSource{key: "int", value: []string{"-1"}},
		&FakeSource{key: "int-bad", value: []string{"not an integer"}},
	}

	this.reader = NewReader(this.sources...)
}

////////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestStrings_Found() {
	value := this.reader.Strings("string")

	this.So(value, should.Resemble, []string{"asdf"})
}
func (this *ReaderTestFixture) TestStrings_NotFound() {
	value := this.reader.Strings("blahblah")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestStringsError_Found() {
	value, err := this.reader.StringsError("string")

	this.So(value, should.Resemble, []string{"asdf"})
	this.So(err, should.BeNil)
}
func (this *ReaderTestFixture) TestStringsError_NotFound() {
	value, err := this.reader.StringsError("81")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestStringsPanic_Found() {
	value := this.reader.StringsPanic("string")

	this.So(value, should.Resemble, []string{"asdf"})
}
func (this *ReaderTestFixture) TestStringsPanic_NotFound() {
	this.So(func() { this.reader.StringsPanic("blahblah") }, should.Panic)
}

func (this *ReaderTestFixture) TestStringsFatal_Found() {
	value := this.reader.StringsFatal("string")

	this.So(value, should.Resemble, []string{"asdf"})
}
func (this *ReaderTestFixture) TestStringsFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.StringsFatal("balhaafslk")
	this.So(err, should.NotBeNil)
}

func (this *ReaderTestFixture) TestStringsDefault_Found() {
	value := this.reader.StringsDefault("string", []string{"default"})

	this.So(value, should.Resemble, []string{"asdf"})
}
func (this *ReaderTestFixture) TestStringsDefault_NotFound() {
	value := this.reader.StringsDefault("blahblah", []string{"default"})

	this.So(value, should.Resemble, []string{"default"})
}

//////////////////////////////////////////////////////////////

func (this *ReaderTestFixture) TestIntsError_Found() {
	value, err := this.reader.IntsError("int")

	this.So(value, should.Resemble, []int{42})
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestIntsError_NotFound() {
	value, err := this.reader.IntsError("asdf")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestIntsError_MalformedValue() {
	value, err := this.reader.IntsError("int-bad")

	this.So(value, should.BeNil)
	this.So(err, should.Equal, MalformedValueError)
}

//////////////////////////////////////////////////////////////

type FakeSource struct {
	key   string
	value []string
}

func (this *FakeSource) Name() string {
	return "fake"
}
func (this *FakeSource) Strings(key string) ([]string, error) {
	if key == this.key {
		return this.value, nil
	}
	return nil, KeyNotFoundError
}
