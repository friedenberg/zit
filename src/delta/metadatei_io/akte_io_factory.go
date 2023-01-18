package metadatei_io

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type nopAkteFactory struct{}

func NopAkteFactory() schnittstellen.AkteIOFactory {
	return nopAkteFactory{}
}

func (_ nopAkteFactory) AkteWriter() (sha.WriteCloser, error) {
	return NewNopWriter(), nil
}

func (_ nopAkteFactory) AkteReader(s sha.ShaLike) (sha.ReadCloser, error) {
	return sha.MakeNopReadCloser(io.NopCloser(bytes.NewBuffer(nil))), nil
}

type nopWriter struct {
	hash hash.Hash
}

func NewNopWriter() (w *nopWriter) {
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

func (w *nopWriter) Sha() (s schnittstellen.Sha) {
	s = sha.FromHash(w.hash)

	return
}
