package sha

import (
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ReadCloser interface {
	io.WriterTo
	io.ReadCloser
	Sha() Sha
}

type WriteCloser interface {
	io.ReaderFrom
	io.WriteCloser
	Sha() Sha
}

type readCloser struct {
	tee  io.Reader
	r    io.Reader
	w    io.Writer
	hash hash.Hash
}

func MakeReadCloser(r io.Reader) (src readCloser) {
	switch rt := r.(type) {
	case *readCloser:
		src = *rt

	case readCloser:
		src = rt

	default:
		src.hash = sha256.New()
		src.r = r
	}

	src.setupTee()

	return
}

func MakeReadCloserTee(r io.Reader, w io.Writer) (src readCloser) {
	switch rt := r.(type) {
	case *readCloser:
		src = *rt
		src.w = w

	case readCloser:
		src = rt
		src.w = w

	default:
		src.hash = sha256.New()
		src.r = r
		src.w = w
	}

	src.setupTee()

	return
}

func (src *readCloser) setupTee() {
	if src.w == nil {
		src.tee = io.TeeReader(src.r, src.hash)
	} else {
		src.tee = io.TeeReader(src.r, io.MultiWriter(src.hash, src.w))
	}
}

func (r readCloser) WriteTo(w io.Writer) (n int64, err error) {
	//TODO-P3 determine why something in the copy returns an EOF
	if n, err = io.Copy(w, r.tee); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (r readCloser) Read(b []byte) (n int, err error) {
	return r.tee.Read(b)
}

func (r readCloser) Close() (err error) {
	if c, ok := r.r.(io.Closer); ok {
		defer errors.DeferredCloser(&err, c)
	}

	if c, ok := r.w.(io.Closer); ok {
		defer errors.DeferredCloser(&err, c)
	}

	return
}

func (r readCloser) Sha() Sha {
	return FromHash(r.hash)
}