package age_io

import (
	"compress/gzip"
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type Reader interface {
	sha.ReadCloser
}

type reader struct {
	hash hash.Hash
	rAge io.Reader
	rZip io.ReadCloser
	tee  io.Reader
}

func NewReader(o ReadOptions) (r *reader, err error) {
	r = &reader{}

	if r.rAge, err = o.Age.Decrypt(o.Reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	r.hash = sha256.New()

	if o.UseZip {
		if r.rZip, err = gzip.NewReader(r.rAge); err != nil {
			err = errors.Wrap(err)
			return
		}

		r.tee = io.TeeReader(r.rZip, r.hash)
	} else {
		r.tee = io.TeeReader(r.rAge, r.hash)
	}

	return
}

func (r *reader) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, r.tee)
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.tee.Read(p)
}

func (r *reader) Close() (err error) {
	if r.rZip == nil {
		return
	}

	if err = r.rZip.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (r *reader) Sha() (s schnittstellen.Sha) {
	s = sha.FromHash(r.hash)

	return
}
