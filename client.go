package chunk

import (
	"fmt"
	"io"
)

type Client struct {
	channel io.ReadWriter
}

func (c *Client) Get(name string) (*Chunk, error) {
	_, err := fmt.Fprintf(c.channel, "%s %s\n", VerbGET, name)
	if err != nil {
		return nil, err
	}
	response, err := consumeToNewline(c.channel)
	if err != nil {
		return nil, err
	}
	if response == ResponseERR {
		return nil, ErrChunkDoesNotExist
	}
	if response != ResponseOK {
		return nil, ErrBadProtocol
	}
	chunk := new(Chunk)
	_, err = chunk.ReadFrom(c.channel)
	if err != nil {
		return chunk, err
	}
	response, err = consumeToNewline(c.channel)
	if response != ResponseDONE {
		return chunk, ErrBadProtocol
	}
	return chunk, err
}

func (c *Client) Set(name string, chunk *Chunk) error {
	_, err := fmt.Fprintf(c.channel, "%s %s\n", VerbSET, name)
	if err != nil {
		return err
	}
	response, err := consumeToNewline(c.channel)
	if err != nil {
		return err
	}
	if response == ResponseERR {
		return ErrChunkExists
	}
	if response != ResponseOK {
		return ErrBadProtocol
	}
	_, err = chunk.WriteTo(c.channel)
	response, err = consumeToNewline(c.channel)
	if response != ResponseDONE {
		return ErrBadProtocol
	}
	return err
}

func (c *Client) Del(name string) error {
	_, err := fmt.Fprintf(c.channel, "%s %s\n", VerbDEL, name)
	if err != nil {
		return err
	}
	response, err := consumeToNewline(c.channel)
	if err != nil {
		return err
	}
	if response == ResponseERR {
		return ErrChunkDoesNotExist
	}
	if response != ResponseOK {
		return ErrBadProtocol
	}
	return nil
}
