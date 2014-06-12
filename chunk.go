package chunk

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type bodyFormat uint16

type Chunk struct {
	Length       uint64
	BodyChecksum uint32
	HeadChecksum uint32
	Body         []byte
}

const (
	headerSize = 16
)

var (
	ErrChecksumMismatch = errors.New("checksum mismatch")
)

func (c *Chunk) ReadFrom(r io.Reader) (n int64, err error) {
	// Get full header.
	buf := make([]byte, headerSize)
	ni, err := io.ReadFull(r, buf)
	n += int64(ni)
	if err != nil {
		return
	}
	if n != headerSize {
		return n, io.ErrUnexpectedEOF
	}
	// Checksum the header.
	headerCRC := crc32.ChecksumIEEE(buf[:12])

	// Read Header
	c.Length = binary.LittleEndian.Uint64(buf[:8])
	c.BodyChecksum = binary.LittleEndian.Uint32(buf[8:12])
	c.HeadChecksum = binary.LittleEndian.Uint32(buf[12:16])

	if headerCRC != c.HeadChecksum {
		return n, ErrChecksumMismatch
	}
	c.Body = make([]byte, c.Length)
	ni, err = io.ReadFull(r, c.Body)
	n += int64(ni)
	if err != nil {
		return
	}
	if bodyChecksum := crc32.ChecksumIEEE(c.Body); bodyChecksum != c.BodyChecksum {
		return n, ErrChecksumMismatch
	}
	return
}

func (c *Chunk) WriteTo(w io.Writer) (p int64, err error) {
	c.Length = uint64(len(c.Body))
	buf := make([]byte, headerSize)
	binary.LittleEndian.PutUint64(buf[0:8], c.Length)
	c.BodyChecksum = crc32.ChecksumIEEE(c.Body)
	binary.LittleEndian.PutUint32(buf[8:12], c.BodyChecksum)
	c.HeadChecksum = crc32.ChecksumIEEE(buf[:12])
	binary.LittleEndian.PutUint32(buf[12:16], c.HeadChecksum)
	pi, err := w.Write(buf)
	p += int64(pi)
	if err != nil {
		return p, err
	}
	pi, err = w.Write(c.Body)
	p += int64(pi)
	return
}
