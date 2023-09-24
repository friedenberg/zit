package standort

import (
	"crypto/sha256"
	"hash"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Reader interface {
	sha.ReadCloser
}

type reader struct {
	hash    hash.Hash
	rAge    io.Reader
	rExpand io.ReadCloser
	tee     io.Reader
}

func NewReader(o ReadOptions) (r *reader, err error) {
	r = &reader{}

	if r.rAge, err = o.Age.Decrypt(o.Reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	if r.rExpand, err = o.CompressionType.NewReader(r.rAge); err != nil {
		err = errors.Wrap(err)
		return
	}

	r.hash = sha256.New()
	r.tee = io.TeeReader(r.rExpand, r.hash)

	return
}

func (r *reader) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, r.tee)
}

func (r *reader) Read(p []byte) (n int, err error) {
	return r.tee.Read(p)
}

func (r *reader) Close() (err error) {
	if err = r.rExpand.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (r *reader) GetShaLike() (s schnittstellen.ShaLike) {
	s = sha.FromHash(r.hash)

	return
}
