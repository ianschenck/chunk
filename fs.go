package chunk

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type FSStore struct {
	BasePath string
}

func (f *FSStore) Get(name string) (io.ReadCloser, error) {
	return os.OpenFile(path.Join(f.BasePath, name), os.O_RDONLY, 0600)
}

func (f *FSStore) Set(name string) (io.WriteCloser, error) {
	return os.OpenFile(path.Join(f.BasePath, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
}

func (f *FSStore) Del(name string) error {
	return os.Remove(path.Join(f.BasePath, name))
}

func (f *FSStore) List(prefix string) []string {
	files, err := ioutil.ReadDir(f.BasePath)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(files))
	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) {
			names = append(names, f.Name())
		}
	}
	return names
}
