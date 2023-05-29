package store_util

import (
	"bytes"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type verzeichnisseElement interface {
	DidRead() bool
	HasChanges() bool
	io.ReaderFrom
	io.WriterTo
}

type verzeichnisseWrapper[T verzeichnisseElement] struct {
	path  string
	index T
}

func makeVerzeichnisseWrapper[T verzeichnisseElement](
	e T,
	path string,
) verzeichnisseWrapper[T] {
	return verzeichnisseWrapper[T]{
		path:  path,
		index: e,
	}
}

func (ei *verzeichnisseWrapper[T]) ReadIfNecessary(
	vf schnittstellen.VerzeichnisseFactory,
) (err error) {
	if ei.index.DidRead() {
		return
	}

	var rc io.ReadCloser

	if rc, err = vf.ReadCloserVerzeichnisse(ei.path); err != nil {
		if errors.IsNotExist(err) {
			err = nil
			rc = sha.MakeReadCloser(bytes.NewBuffer(nil))
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, rc)

	if _, err = ei.index.ReadFrom(rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ei *verzeichnisseWrapper[T]) Get(
	vf schnittstellen.VerzeichnisseFactory,
) (i T, err error) {
	if err = ei.ReadIfNecessary(vf); err != nil {
		err = errors.Wrap(err)
		return
	}

	i = ei.index

	return
}

func (ei *verzeichnisseWrapper[T]) Flush(
	vf schnittstellen.VerzeichnisseFactory,
) (err error) {
	if !ei.index.HasChanges() {
		return
	}

	var wc schnittstellen.ShaWriteCloser

	if wc, err = vf.WriteCloserVerzeichnisse(ei.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, wc)

	if _, err = ei.index.WriteTo(wc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
