package chunk

import (
	"bytes"
	"testing"
)

func TestChunkReadFrom(t *testing.T) {
	c := &Chunk{}

	buf := []byte{
		12, 0, 0, 0, 0, 0, 0, 0, // Length
		163, 28, 41, 28, // BodyChecksum
		181, 164, 122, 100, // HeadChecksum
		72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33,
	}

	b := bytes.NewBuffer(buf)

	_, err := c.ReadFrom(b)
	if err != nil {
		t.Fatal(err)
	}

	if string(c.Body) != "Hello World!" {
		t.Error(string(c.Body))
	}
}

func TestChunkWriteTo(t *testing.T) {
	c := &Chunk{}

	b := &bytes.Buffer{}

	c.Body = []byte("Hello World!")

	_, err := c.WriteTo(b)

	if err != nil {
		t.Error(err)
	}

	// Expected output:
	buf := []byte{
		12, 0, 0, 0, 0, 0, 0, 0, // Length
		163, 28, 41, 28, // BodyChecksum
		181, 164, 122, 100, // HeadChecksum
		72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33,
	}

	ba := b.Bytes()
	for i := range ba {
		if ba[i] != buf[i] {
			t.Fatalf("\n%v\n%v\n", ba, buf)
		}
	}
}

func TestCopy(t *testing.T) {
	src := &Chunk{}
	dst := &Chunk{}

	src.Body = []byte("Hello World!")

	b := &bytes.Buffer{}

	p, err := src.WriteTo(b)
	if err != nil {
		t.Fatal(err)
	}
	n, err := dst.ReadFrom(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != p {
		t.Fatal("read bytes does not equal write bytes")
	}
	if string(src.Body) != string(dst.Body) {
		t.Fatalf("'%s' != '%s'", string(src.Body), string(dst.Body))
	}
}
