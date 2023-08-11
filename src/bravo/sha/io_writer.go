package sha

import (
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type writer struct {
	c io.Closer
	w io.Writer
	h hash.Hash
}

func MakeWriter(in io.Writer) (out *writer) {
	h := sha256.New()

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
	if w.c == nil {
		return
	}

	if err = w.c.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *writer) GetShaLike() (s schnittstellen.ShaLike) {
	s = FromHash(w.h)

	return
}
