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
		&FakeSource{key: "string-no-values", value: []string{}},
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

func (this *ReaderTestFixture) TestString_Found() {
	value := this.reader.String("string")

	this.So(value, should.Equal, "asdf")
}
func (this *ReaderTestFixture) TestString_NotFound() {
	value := this.reader.String("blahblah")

	this.So(value, should.BeEmpty)
}

func (this *ReaderTestFixture) TestStringError_Found() {
	value, err := this.reader.StringError("string")

	this.So(value, should.Resemble, "asdf")
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestStringError_NotFound() {
	value, err := this.reader.StringError("81")

	this.So(value, should.Equal, "")
	this.So(err, should.Equal, KeyNotFoundError)
}

func (this *ReaderTestFixture) TestStringError_FoundButNoValuesProvided() {
	value, err := this.reader.StringError("string-no-values")

	this.So(value, should.Equal, "")
	this.So(err, should.BeNil)
}

func (this *ReaderTestFixture) TestStringPanic_Found() {
	value := this.reader.StringPanic("string")

	this.So(value, should.Resemble, "asdf")
}

func (this *ReaderTestFixture) TestStringPanic_NotFound() {
	this.So(func() { this.reader.StringPanic("blahblah") }, should.Panic)
}

func (this *ReaderTestFixture) TestStringFatal_Found() {
	value := this.reader.StringFatal("string")

	this.So(value, should.Resemble, "asdf")
}

func (this *ReaderTestFixture) TestStringFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.StringFatal("balhaafslk")
	this.So(err, should.NotBeNil)
}

func (this *ReaderTestFixture) TestStringDefault_Found() {
	value := this.reader.StringDefault("string", "default")

	this.So(value, should.Resemble, "asdf")
}

func (this *ReaderTestFixture) TestStringDefault_NotFound() {
	value := this.reader.StringDefault("blahblah", "default")

	this.So(value, should.Resemble, "default")
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

func (this *ReaderTestFixture) TestInts_Found() {
	value := this.reader.Ints("int")

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestInts_NotFound() {
	value := this.reader.Ints("qrew")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestInts_MalformedValue() {
	value := this.reader.Ints("int-bad")

	this.So(value, should.BeNil)
}

func (this *ReaderTestFixture) TestIntsPanic_Found() {
	value := this.reader.IntsPanic("int")

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestIntsPanic_NotFound() {
	this.So(func() { this.reader.IntsPanic("blah blah") }, should.Panic)
}

func (this *ReaderTestFixture) TestIntsPanic_MalformedValue() {
	this.So(func() { this.reader.IntsPanic("int-bad") }, should.Panic)
}

func (this *ReaderTestFixture) TestIntsFatal_Found() {
	value := this.reader.IntsFatal("int")

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestIntsFatal_NotFound() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.IntsFatal("balhaafslk")
	this.So(err, should.NotBeNil)
}

func (this *ReaderTestFixture) TestIntsFatal_MalformedValue() {
	var err error
	fatal = func(e error) { err = e }
	this.reader.IntsFatal("int-bad")
	this.So(err, should.NotBeNil)
}

func (this *ReaderTestFixture) TestIntsDefault_Found() {
	value := this.reader.IntsDefault("int", []int{84})

	this.So(value, should.Resemble, []int{42})
}

func (this *ReaderTestFixture) TestIntsDefault_NotFound() {
	value := this.reader.IntsDefault("missing", []int{84})

	this.So(value, should.Resemble, []int{84})
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

//////////////////////////////////////////////////////////////
