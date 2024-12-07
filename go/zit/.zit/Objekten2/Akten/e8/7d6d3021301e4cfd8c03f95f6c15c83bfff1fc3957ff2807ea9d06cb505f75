package sha

import (
	"hash"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type writer struct {
	closed bool
	c      io.Closer
	w      io.Writer
	h      hash.Hash
	sh     Sha
}

func MakeWriter(in io.Writer) (out *writer) {
	h := hash256Pool.Get()

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
	w.setShaLikeIfNecessary()

	if w.c == nil {
		return
	}

	if err = w.c.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *writer) setShaLikeIfNecessary() {
	if w.h != nil {
		errors.PanicIfError(w.sh.SetFromHash(w.h))

		hash256Pool.Put(w.h)
		w.h = nil
	}
}

func (w *writer) GetShaLike() (s interfaces.Sha) {
	w.setShaLikeIfNecessary()
	s = &w.sh

	return
}
