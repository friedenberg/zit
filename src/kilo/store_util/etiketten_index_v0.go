package store_util

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/etiketten_index"
)

type etikettenIndexV0 struct {
	lock    *sync.Mutex
	didRead bool
	index   etiketten_index.Index
}

func makeEtikettenIndexV0() *etikettenIndexV0 {
	return &etikettenIndexV0{
		lock:  &sync.Mutex{},
		index: etiketten_index.MakeIndex(),
	}
}

func (ei *etikettenIndexV0) ReadIfNecessary(
	vf schnittstellen.VerzeichnisseFactory,
) (err error) {
	ei.lock.Lock()
	defer ei.lock.Unlock()

	if ei.didRead {
		return
	}

	var rc schnittstellen.ShaReadCloser

	if rc, err = vf.ReadCloserVerzeichnisse("EtikettenIndexV0"); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	if _, err = ei.index.ReadFrom(rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ei *etikettenIndexV0) GetEtikettenIndex(
	vf schnittstellen.VerzeichnisseFactory,
) (i etiketten_index.Index, err error) {
	if err = ei.ReadIfNecessary(vf); err != nil {
		err = errors.Wrap(err)
		return
	}

	i = ei.index

	return
}

func (ei *etikettenIndexV0) Flush(
	vf schnittstellen.VerzeichnisseFactory,
) (err error) {
	var wc schnittstellen.ShaWriteCloser

	if wc, err = vf.WriteCloserVerzeichnisse("EtikettenIndexV0"); err != nil {
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
