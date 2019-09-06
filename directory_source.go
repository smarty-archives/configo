package configo

import (
	"io/ioutil"
	"log"
	"strings"
)

// DirectorySource makes a key-value mapping available of filename (key) to file contents (value).
type DirectorySource struct {
	files     map[string]string
	mustExist bool
	path      string
}

// FromDirectory reads the directory path provided. If the path does not exist, a panic will result.
func FromDirectory(path string) *DirectorySource {
	return &DirectorySource{
		mustExist: true,
		path:      path,
	}
}

func FromOptionalDirectories(paths ...string) MultiSource {
	var sources MultiSource
	for _, path := range paths {
		sources = append(sources, FromOptionalDirectory(path))
	}
	return sources
}

// FromOptionalDirectory reads the directory path provided, if it exists.
func FromOptionalDirectory(path string) *DirectorySource {
	return &DirectorySource{
		mustExist: false,
		path:      path,
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

	return nil, ErrKeyNotFound
}

func (this *DirectorySource) Initialize() {
	this.files = make(map[string]string, 32)

	if files, err := ioutil.ReadDir(this.path); err != nil {
		log.Printf("[INFO] directory not read [%s]: %s\n", this.path, err)
		if this.mustExist {
			panic("directory must exist")
		}
	} else {
		for _, file := range files {
			if !file.IsDir() {
				key := sanitizeKey(strings.ToLower(file.Name()))
				this.files[key] = file.Name()
			}
		}
	}
}
