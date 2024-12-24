package dir_layout

import (
	"crypto/sha256"
	"hash"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
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

func (r *reader) Seek(offset int64, whence int) (actual int64, err error) {
	seeker, ok := r.rAge.(io.Seeker)

	if !ok {
		err = errors.Errorf("seeking not supported")
		return
	}

	return seeker.Seek(offset, whence)
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

func (r *reader) GetShaLike() (s interfaces.Sha) {
	s = sha.FromHash(r.hash)

	return
}
