package objekte

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
)

type Writer interface {
	io.WriteCloser
	Sha() sha.Sha
}

type writer struct {
	hash       hash.Hash
	tee        io.Writer
	wZip, wAge io.WriteCloser
	wBuf       *bufio.Writer
}

func NewZippedWriter(age age.Age, out io.Writer) (w *writer, err error) {
	if w, err = NewWriter(age, out); err != nil {
		err = errors.Error(err)
		return
	}

	w.wZip = gzip.NewWriter(w.wAge)
	w.tee = io.MultiWriter(w.hash, w.wZip)

	return
}

func NewWriter(age age.Age, out io.Writer) (w *writer, err error) {
	w = &writer{}

	w.wBuf = bufio.NewWriter(out)

	if w.wAge, err = age.Encrypt(out); err != nil {
		err = errors.Error(err)
		return
	}

	w.hash = sha256.New()

	w.tee = io.MultiWriter(w.hash, w.wAge)

	return
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.tee.Write(p)
}

func (w *writer) Close() (err error) {
	if w.wZip != nil {
		if err = w.wZip.Close(); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if err = w.wAge.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (w *writer) Sha() (s sha.Sha) {
	s = sha.FromHash(w.hash)

	return
}
