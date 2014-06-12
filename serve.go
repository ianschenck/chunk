package chunk

import (
	"io"
)

func Serve(store Store, rw io.ReadWriter) (err error) {
	verb := make([]byte, 3)
	_, err = rw.Read(verb)
	if err != nil {
		return
	}
	switch string(verb) {
	case VerbGET:
		if err := doGet(store, rw); err != nil {
			return err
		}
	case VerbDEL:
		if err := doDel(store, rw); err != nil {
			return err
		}
	case VerbSET:
		if err := doSet(store, rw); err != nil {
			return err
		}
	default:
		return ErrBadProtocol
	}
	return nil
}

func doGet(store Store, rw io.ReadWriter) (err error) {
	// consume space
	if !consumeSpace(rw) {
		return ErrBadProtocol
	}
	// read chunk name
	name, err := consumeToNewline(rw)
	if err != nil {
		return err
	}
	if name == "" {
		return ErrBadProtocol
	}
	r, err := store.Get(name)
	if err != nil {
		_, err = writeErr(rw)
		if err != nil {
			return err
		}
		return nil
	}
	defer r.Close()
	if _, err = writeOk(rw); err != nil {
		return err
	}
	if _, err = io.Copy(rw, r); err != nil {
		return
	}
	_, err = writeDone(rw)
	return
}

func doDel(store Store, rw io.ReadWriter) (err error) {
	// consume space
	if !consumeSpace(rw) {
		return ErrBadProtocol
	}
	// read chunk name
	name, err := consumeToNewline(rw)
	if err != nil {
		return err
	}
	if name == "" {
		return ErrBadProtocol
	}
	if err = store.Del(name); err != nil {
		_, err = writeErr(rw)
		return err
	}
	_, err = writeOk(rw)
	return
}

func doSet(store Store, rw io.ReadWriter) (err error) {
	// consume space
	if !consumeSpace(rw) {
		return ErrBadProtocol
	}
	// read chunk name
	name, err := consumeToNewline(rw)
	if err != nil {
		return err
	}
	if name == "" {
		return ErrBadProtocol
	}
	w, err := store.Set(name)
	if err != nil {
		_, err = writeErr(rw)
		if err != nil {
			return err
		}
		return nil
	}
	defer w.Close()
	if _, err = writeOk(rw); err != nil {
		return err
	}
	chunk := new(Chunk)
	_, err = chunk.ReadFrom(rw)
	if err != nil {
		_, err = writeErr(rw)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = chunk.WriteTo(w)
	if err != nil {
		return
	}
	_, err = writeDone(rw)
	return
}
