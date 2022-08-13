package io_builder

import (
	"io"

	"github.com/friedenberg/zit/bravo/errors"
)

type ReaderMakeFunc func() (io.ReadCloser, error)
type ReaderChainFunc func(io.Reader) (io.ReadCloser, error)

type closerNull struct {
	io.Reader
}

func (c closerNull) Close() error {
	return nil
}

func MakeReaderMakeFunc(f func() (io.Reader, error)) ReaderMakeFunc {
	return func() (rc io.ReadCloser, err error) {
		var r io.Reader

		if r, err = f(); err != nil {
			err = errors.Error(err)
			return
		}

		rc = closerNull{Reader: r}

		return
	}
}

func Reader(mf ReaderMakeFunc, rcs ...ReaderChainFunc) (r io.ReadCloser, err error) {
	if r, err = mf(); err != nil {
		err = errors.Wrapped(err, "failed to make reader")
		return
	}

	for _, rc := range rcs {
		var nr io.ReadCloser

		if nr, err = rc(r); err != nil {
			err = errors.Wrapped(err, "failed to chain reader")
			return
		}

		r = nr
	}

	return
}
