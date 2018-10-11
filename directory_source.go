package configo

import (
	"io/ioutil"
	"strings"
)

// DirectorySource makes a key-value mapping available of filename (key) to file contents (value).
type DirectorySource struct {
	files map[string]string
	mustExist bool
	path string
}


// Reads the directory path provided. If the path does not exist, a panic will result.
func FromDirectory(path string) *DirectorySource {
	return &DirectorySource{
		mustExist: true,
		path: path,
	}
}


// Reads the directory path provided, if it exists.
func FromOptionalDirectory(path string) *DirectorySource {
	return &DirectorySource{
		mustExist: false,
		path: path,
	}
}


func (this *DirectorySource) Strings(key string) ([]string, error) {
	key = sanitizeKey(strings.ToLower(key))

	if filename, found := this.files[key]; found {
		path := this.path + "/" + filename
		if data, err := ioutil.ReadFile(path); err == nil {
			return []string{string(data)}, nil
		} else {
			return nil, err
		}
	}

	return nil, KeyNotFoundError
}


func (this *DirectorySource) Initialize() {
	this.files = make(map[string]string, 32)

	if files, err := ioutil.ReadDir(this.path); err != nil {
		if this.mustExist {
			panic(err)
		}
	} else {
		for _, file := range files {
			key := sanitizeKey(strings.ToLower(file.Name()))
			this.files[key] = file.Name()
		}
	}
}
