package sha

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

// TODO-P4 remove
type (
	ReadCloser  = schnittstellen.ShaReadCloser
	WriteCloser = schnittstellen.ShaWriteCloser
)

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
	// TODO-P3 determine why something in the copy returns an EOF
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

func (r readCloser) GetShaLike() schnittstellen.ShaLike {
	return FromHash(r.hash)
}

type nopReadCloser struct {
	io.ReadCloser
}

func MakeNopReadCloser(rc io.ReadCloser) ReadCloser {
	return nopReadCloser{
		ReadCloser: rc,
	}
}

func (nrc nopReadCloser) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, nrc.ReadCloser)
}

func (nrc nopReadCloser) GetShaLike() schnittstellen.ShaLike {
	return Sha{}
}

type nopAkteFactory struct{}

func NopAkteFactory() schnittstellen.AkteIOFactory {
	return nopAkteFactory{}
}

func (_ nopAkteFactory) AkteWriter() (WriteCloser, error) {
	return MakeNopWriter(), nil
}

func (_ nopAkteFactory) AkteReader(s ShaLike) (ReadCloser, error) {
	return MakeNopReadCloser(io.NopCloser(bytes.NewBuffer(nil))), nil
}

type nopWriter struct {
	hash hash.Hash
}

func MakeNopWriter() (w *nopWriter) {
	w = &nopWriter{
		hash: sha256.New(),
	}

	return
}

func (w *nopWriter) ReadFrom(r io.Reader) (n int64, err error) {
	if n, err = io.Copy(w.hash, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *nopWriter) Write(p []byte) (n int, err error) {
	return w.hash.Write(p)
}

func (w *nopWriter) Close() (err error) {
	return
}

func (w *nopWriter) GetShaLike() (s schnittstellen.ShaLike) {
	s = FromHash(w.hash)

	return
}
