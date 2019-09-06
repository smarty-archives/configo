package configo

import (
	"io/ioutil"
	"log"
	"path"
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

func FromOptionalDirectories(directories ...string) MultiSource {
	var sources MultiSource
	for _, directory := range directories {
		sources = append(sources, FromOptionalDirectory(directory))
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

	filename, found := this.files[key]
	if !found {
		return nil, ErrKeyNotFound
	}

	data, err := ioutil.ReadFile(path.Join(this.path, filename))
	if err != nil {
		return nil, err
	}

	return []string{string(data)}, nil
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
