package configo

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestDirectorySourceFixture(t *testing.T) {
	gunit.Run(new(DirectorySourceFixture), t)
}

type DirectorySourceFixture struct {
	*gunit.Fixture

	dirPath string
	files   map[string]string
}

////////////////////////////////////////////////////////////////

// test c'tor
func (this *DirectorySourceFixture) Setup() {
	if path, err := ioutil.TempDir("", "dirSrc"); err == nil {
		this.dirPath = path
		this.files = map[string]string{
			"File1":                   "My file contents",
			"A-name_withMixed&Casing": "ContEnT$ _with_\n{{mixed}} Casing\n",
		}

		for filename, content := range this.files {
			full := path + "/" + filename
			if err := ioutil.WriteFile(full, []byte(content), 0600); err != nil {
				this.Error("Error writing test files:", err)
			}
		}
	} else {
		panic(err)
	}
}

// test d'tor
func (this *DirectorySourceFixture) Teardown() {
	if err := os.RemoveAll(this.dirPath); err != nil {
		this.Error("Error removing test files:", err)
	}
}

////////////////////////////////////////////////////////////////

func (this *DirectorySourceFixture) TestBadDirectoryPanic() {
	src := FromDirectory("&path/@should/!not/*exist")
	this.So(func() { src.Initialize() }, should.Panic)
}

func (this *DirectorySourceFixture) TestBadDirectoryNotPanic() {
	src := FromOptionalDirectory("&path/@should/!not/*exist")
	src.Initialize()
	this.So(len(src.files), should.Equal, 0)
}

func (this *DirectorySourceFixture) TestInitalize() {
	src := FromDirectory(this.dirPath)
	src.Initialize()
	this.So(src.path, should.Equal, this.dirPath)
	this.So(src.mustExist, should.Equal, true)
	this.So(len(src.files), should.Equal, 2)
}

func (this *DirectorySourceFixture) TestStrings() {
	src := FromDirectory(this.dirPath)
	src.Initialize()

	for key, val := range this.files {
		data, err := src.Strings(key)
		this.So(err, should.BeEmpty)
		this.So(len(data), should.NotBeZeroValue)
		this.So(data[0], should.Equal, val)
	}
}

func (this *DirectorySourceFixture) TestStringsCase() {
	src := FromDirectory(this.dirPath)
	src.Initialize()

	for key, val := range this.files {
		key = strings.ToUpper(key)
		data, err := src.Strings(key)
		this.So(err, should.BeEmpty)
		this.So(len(data), should.NotBeZeroValue)
		this.So(data[0], should.Equal, val)
	}
}

func (this *DirectorySourceFixture) TestStringsMissing() {
	src := FromDirectory(this.dirPath)
	src.Initialize()

	data, err := src.Strings("key/does/not-exist")
	this.So(data, should.BeEmpty)
	this.So(err, should.Equal, KeyNotFoundError)
}
