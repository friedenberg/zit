package metadatei_io

import (
	"bytes"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type AkteIOFactory interface {
	AkteReaderFactory
	AkteWriterFactory
}

type AkteReaderFactory interface {
	AkteReader(sha.Sha) (sha.ReadCloser, error)
}

type AkteWriterFactory interface {
	AkteWriter() (sha.WriteCloser, error)
}

type AkteIOFactoryFactory interface {
	AkteFactory(gattung.Gattung) AkteIOFactory
}

type nopAkteFactory struct{}

func NopAkteFactory() AkteIOFactory {
	return nopAkteFactory{}
}

func (_ nopAkteFactory) AkteWriter() (sha.WriteCloser, error) {
	return NewNopWriter(), nil
}

func (_ nopAkteFactory) AkteReader(s sha.Sha) (sha.ReadCloser, error) {
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

func (w *nopWriter) Sha() (s sha.Sha) {
	s = sha.FromHash(w.hash)

	return
}
