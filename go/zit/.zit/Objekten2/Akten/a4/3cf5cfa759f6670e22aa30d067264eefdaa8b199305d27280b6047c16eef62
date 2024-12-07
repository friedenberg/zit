package dir_layout

import (
	"bufio"
	"crypto/sha256"
	"hash"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type Writer interface {
	sha.WriteCloser
}

type writer struct {
	hash            hash.Hash
	tee             io.Writer
	wCompress, wAge io.WriteCloser
	wBuf            *bufio.Writer
}

func NewWriter(o WriteOptions) (w *writer, err error) {
	w = &writer{}

	w.wBuf = bufio.NewWriter(o.Writer)

	if w.wAge, err = o.Encrypt(w.wBuf); err != nil {
		err = errors.Wrap(err)
		return
	}

	w.hash = sha256.New()

	w.wCompress = o.CompressionType.NewWriter(w.wAge)
	w.tee = io.MultiWriter(w.hash, w.wCompress)

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
	if err = w.wCompress.Close(); err != nil {
		err = errors.Wrap(err)
		return
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

func (w *writer) GetShaLike() (s interfaces.Sha) {
	s = sha.FromHash(w.hash)

	return
}
