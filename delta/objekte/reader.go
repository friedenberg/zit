package objekte

import (
	"compress/gzip"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/age"
)

type Reader interface {
	io.ReadCloser
	Sha() sha.Sha
}

type reader struct {
	hash hash.Hash
	rAge io.Reader
	rZip io.ReadCloser
	tee  io.Reader
}

func NewZippedReader(age age.Age, in io.Reader) (r *reader, err error) {
	if r, err = NewReader(age, in); err != nil {
		err = errors.Error(err)
		return
	}

	if r.rZip, err = gzip.NewReader(r.rAge); err != nil {
		err = errors.Error(err)
		return
	}

	r.tee = io.TeeReader(r.rZip, r.hash)

	return
}

func NewReader(age age.Age, in io.Reader) (r *reader, err error) {
	r = &reader{}

	if r.rAge, err = age.Decrypt(in); err != nil {
		err = errors.Error(err)
		return
	}

	r.hash = sha256.New()

	r.tee = io.TeeReader(r.rAge, r.hash)

	return
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.tee.Read(p)
}

func (r *reader) Close() (err error) {
	if r.rZip == nil {
		return
	}

	if err = r.rZip.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (r *reader) Sha() (s sha.Sha) {
	s = sha.FromHash(r.hash)

	return
}
