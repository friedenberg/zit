package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type Page struct {
	lock *sync.Mutex
	pageId
	ioFactory
	pool        *collections.Pool[zettel.Transacted, *zettel.Transacted]
	added       zettel.HeapTransacted
	flushFilter schnittstellen.FuncIter[*zettel.Transacted]
	State
}

func makeZettelenPage(
	iof ioFactory,
	pid pageId,
	pool *collections.Pool[zettel.Transacted, *zettel.Transacted],
	fff ZettelTransactedWriterGetter,
) (p *Page) {
	flushFilter := collections.MakeWriterNoop[*zettel.Transacted]()

	if fff != nil {
		flushFilter = fff.ZettelTransactedWriter(pid.index)
	}

	p = &Page{
		lock:        &sync.Mutex{},
		ioFactory:   iof,
		pageId:      pid,
		pool:        pool,
		added:       zettel.MakeHeapTransacted(),
		flushFilter: flushFilter,
	}

	p.added.SetPool(p.pool)

	return
}

func (zp *Page) Add(z *zettel.Transacted) (err error) {
	if z == nil {
		err = errors.Errorf("trying to add nil zettel")
		return
	}

	if err = zp.flushFilter(z); err != nil {
		if collections.IsStopIteration(err) {
			errors.Log().Printf("eliding %s", z.Kennung())
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	zp.lock.Lock()
	defer zp.lock.Unlock()

	zp.added.Add(*z)
	zp.State = StateChanged

	return
}

func (zp *Page) Flush() (err error) {
	zp.lock.Lock()
	defer zp.lock.Unlock()

	if zp.State < StateChanged {
		return
	}

	errors.Log().Printf("flushing page: %s", zp.path)

	var w io.WriteCloser

	if w, err = zp.WriteCloserVerzeichnisse(zp.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	w1 := bufio.NewWriter(w)

	defer errors.DeferredFlusher(&err, w1)

	if _, err = zp.writeTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.added.Reset()

	return
}

func (zp *Page) copy(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (n int64, err error) {
	var r1 io.ReadCloser

	if r1, err = zp.ReadCloserVerzeichnisse(zp.path); err != nil {
		if errors.IsNotExist(err) {
			r1 = ioutil.NopCloser(bytes.NewReader(nil))
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, r1)

	r := bufio.NewReader(r1)

	dec := gob.NewDecoder(r)

	if err = zp.added.MergeStream(
		func() (tz *zettel.Transacted, err error) {
			tz = zp.pool.Get()

			if err = dec.Decode(tz); err != nil {
				if errors.IsEOF(err) {
					err = collections.MakeErrStopIteration()
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
		w,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *Page) writeTo(w1 io.Writer) (n int64, err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	if n, err = zp.copy(
		collections.MakeChain(
			zp.flushFilter,
			zettel.MakeWriterGobEncoder(w).WriteZettelVerzeichnisse,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *Page) Copy(
	w schnittstellen.FuncIter[*zettel.Transacted],
) (n int64, err error) {
	acquired := zp.lock.TryLock()

	if !acquired {
		err = MakeErrConcurrentPageAccess()
		return
	}

	defer zp.lock.Unlock()

	return zp.copy(w)
}

func (zp *Page) WriteTo(w1 io.Writer) (n int64, err error) {
	acquired := zp.lock.TryLock()

	if !acquired {
		err = MakeErrConcurrentPageAccess()
		return
	}

	defer zp.lock.Unlock()

	return zp.writeTo(w1)
}
