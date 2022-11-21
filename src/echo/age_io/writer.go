package age_io

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Writer interface {
	sha.WriteCloser
}

type writer struct {
	hash       hash.Hash
	tee        io.Writer
	wZip, wAge io.WriteCloser
	wBuf       *bufio.Writer
}

func NewWriter(o WriteOptions) (w *writer, err error) {
	w = &writer{}

	w.wBuf = bufio.NewWriter(o.Writer)

	if w.wAge, err = o.Age.Encrypt(o.Writer); err != nil {
		err = errors.Wrap(err)
		return
	}

	w.hash = sha256.New()

	w.tee = io.MultiWriter(w.hash, w.wAge)

	if o.UseZip {
		w.wZip = gzip.NewWriter(w.wAge)
		w.tee = io.MultiWriter(w.hash, w.wZip)
	}

	return
}

func (w *writer) ReadFrom(r io.Reader) (n int64, err error) {
	if n, err = io.Copy(w.tee, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.tee.Write(p)
}

func (w *writer) Close() (err error) {
	if w.wZip != nil {
		if err = w.wZip.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = w.wAge.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *writer) Sha() (s sha.Sha) {
	s = sha.FromHash(w.hash)

	return
}
