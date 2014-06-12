package chunk

import (
	"io"
)

type Store interface {
	Get(name string) (io.ReadCloser, error)
	Set(name string) (io.WriteCloser, error)
	Del(name string) error
	List(prefix string) []string
}
