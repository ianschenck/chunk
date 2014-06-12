// Package chunk defines the Chunk storage primitive and protocol.
package chunk

import (
	"errors"
	"io"
)

var (
	// ErrBadProtocol is a malformed interaction.
	ErrBadProtocol = errors.New("bad protocol")

	ErrChunkExists       = errors.New("chunk already exists")
	ErrChunkDoesNotExist = errors.New("chunk does not exist")
)

const (
	newline = 10
	space   = 32

	VerbGET = "GET"
	VerbSET = "SET"
	VerbDEL = "DEL"

	ResponseOK   = "OK!"
	ResponseERR  = "ERR"
	ResponseDONE = "DON"
)

func consumeSpace(r io.Reader) bool {
	b := [1]byte{0}
	p, err := r.Read(b[:])
	return err == nil && p == 1 && b[0] == space
}

func consumeToNewline(r io.Reader) (string, error) {
	buf := make([]byte, 0, 1024)
	b := [1]byte{0}
	for {
		p, err := r.Read(b[:])
		if p != 1 || err != nil {
			return string(buf), err
		}
		if b[0] == newline {
			return string(buf), nil
		}
		buf = append(buf, b[0])
	}
}

func writeErr(w io.Writer) (n int, err error) {
	return w.Write([]byte(ResponseERR + "\n"))
}

func writeOk(w io.Writer) (n int, err error) {
	return w.Write([]byte(ResponseOK + "\n"))
}

func writeDone(w io.Writer) (n int, err error) {
	return w.Write([]byte(ResponseDONE + "\n"))
}
