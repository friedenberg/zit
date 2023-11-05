package sha

import (
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type writer struct {
	closed bool
	c      io.Closer
	w      io.Writer
	h      hash.Hash
	sh     Sha
}

func MakeWriter(in io.Writer) (out *writer) {
	h := shaPool.Get()

	if in == nil {
		in = io.Discard
	}

	out = &writer{
		w: io.MultiWriter(h, in),
		h: h,
	}

	if c, ok := in.(io.Closer); ok {
		out.c = c
	}

	return
}

func (w *writer) ReadFrom(r io.Reader) (n int64, err error) {
	if n, err = io.Copy(w.w, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *writer) Close() (err error) {
	w.closed = true

	if w.c == nil {
		return
	}

	if err = w.c.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	w.setShaLikeIfNecessary()

	return
}

func (w *writer) setShaLikeIfNecessary() {
	if !w.closed {
		w.sh = FromHash(w.h)

		shaPool.Put(w.h)
		w.h = nil
	}
}

func (w *writer) GetShaLike() (s schnittstellen.ShaLike) {
	w.setShaLikeIfNecessary()
	s = w.sh

	return
}
