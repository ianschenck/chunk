package chunk

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestChunkGet(t *testing.T) {
	store := MemStore{make(map[string][]byte)}
	lhs, rhs := NewChannel()
	client := Client{lhs}

	go Serve(store, rhs)
	chunk, err := client.Get("foo")
	if err != ErrChunkDoesNotExist {
		t.Errorf("expected '%s', got '%s' for chunk 'foo'", ErrChunkDoesNotExist, err)
	}

	// Create a chunk.
	temp := &bytes.Buffer{}
	chunk = new(Chunk)
	chunk.Body = []byte("Hello World!")
	chunk.WriteTo(temp)
	store.items["foo"] = temp.Bytes()

	go Serve(store, rhs)
	chunk, err = client.Get("foo")
	if err != nil {
		t.Errorf("unexpected error '%s'", err)
	}
}

func TestChunkSet(t *testing.T) {
	store := MemStore{make(map[string][]byte)}
	lhs, rhs := NewChannel()
	client := Client{lhs}

	temp := &bytes.Buffer{}
	chunk := new(Chunk)
	chunk.Body = []byte("Hello World!")
	chunk.WriteTo(temp)

	go Serve(store, rhs)
	err := client.Set("foo", chunk)
	if err != nil {
		t.Errorf("unexpected error '%s'", err)
	}

	go Serve(store, rhs)
	err = client.Set("foo", chunk)
	if err != ErrChunkExists {
		t.Errorf("expected '%s', got '%s' for chunk 'foo'", ErrChunkExists, err)
	}
}

func TestChunkDel(t *testing.T) {
	store := MemStore{make(map[string][]byte)}
	lhs, rhs := NewChannel()
	client := Client{lhs}

	store.items["bar"] = []byte{}

	go Serve(store, rhs)
	err := client.Del("foo")
	if err != ErrChunkDoesNotExist {
		t.Errorf("expected '%s', got '%s' for chunk 'foo'", ErrChunkDoesNotExist, err)
	}

	go Serve(store, rhs)
	err = client.Del("bar")
	if err != nil {
		t.Errorf("unexpected error '%s'", err)
	}

	go Serve(store, rhs)
	err = client.Del("bar")
	if err != ErrChunkDoesNotExist {
		t.Errorf("expected '%s', got '%s' for chunk 'foo'", ErrChunkDoesNotExist, err)
	}
}

type commChannel struct {
	in, out *os.File
}

func (c commChannel) Read(p []byte) (n int, err error) {
	return c.in.Read(p)
}

func (c commChannel) Write(p []byte) (n int, err error) {
	n, err = c.out.Write(p)
	c.out.Sync()
	return
}

func NewChannel() (commChannel, commChannel) {
	lhs_r, lhs_w, _ := os.Pipe()
	rhs_r, rhs_w, _ := os.Pipe()
	return commChannel{lhs_r, rhs_w}, commChannel{rhs_r, lhs_w}
}

type MemStore struct {
	items map[string][]byte
}

func (s MemStore) Get(name string) (io.ReadCloser, error) {
	b, ok := s.items[name]
	if !ok {
		return nil, ErrChunkDoesNotExist
	}
	return ioutil.NopCloser(bytes.NewBuffer(b)), nil
}

func (s MemStore) Set(name string) (io.WriteCloser, error) {
	if _, ok := s.items[name]; ok {
		return nil, ErrChunkExists
	}
	b := &bytes.Buffer{}
	s.items[name] = b.Bytes()
	return noopWriter{b}, nil
}

func (s MemStore) Del(name string) error {
	if _, ok := s.items[name]; !ok {
		return ErrChunkDoesNotExist
	}
	delete(s.items, name)
	return nil
}

func (s MemStore) List(prefix string) []string {
	list := make([]string, 0, len(s.items))
	for item := range s.items {
		list = append(list, item)
	}
	return list
}

type noopWriter struct {
	io.Writer
}

func (n noopWriter) Close() error {
	return nil
}
