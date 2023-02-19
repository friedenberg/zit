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
	pool        schnittstellen.Pool[zettel.Transacted, *zettel.Transacted]
	added       zettel.HeapTransacted
	flushFilter schnittstellen.FuncIter[*zettel.Transacted]
	State
}

func makeZettelenPage(
	iof ioFactory,
	pid pageId,
	pool schnittstellen.Pool[zettel.Transacted, *zettel.Transacted],
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

func (zp *Page) doTryLock() (ok bool) {
	errors.Log().Caller(1, "acquiring lock: %d", zp.pageId.index)
	ok = zp.lock.TryLock()

	if ok {
		errors.Log().Caller(1, "acquired lock: %d", zp.pageId.index)
	} else {
		errors.Log().Caller(1, "failed to acquire lock: %d", zp.pageId.index)
	}

	return
}

func (zp *Page) doLock() {
	errors.Log().Caller(1, "acquiring lock: %d", zp.pageId.index)
	zp.lock.Lock()
	errors.Log().Caller(1, "acquired lock: %d", zp.pageId.index)
}

func (zp *Page) doUnlock() {
	errors.Log().Caller(1, "releasing lock: %d", zp.pageId.index)
	zp.lock.Unlock()
	errors.Log().Caller(1, "released lock: %d", zp.pageId.index)
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

	// do not lock as checking operations perform additions while reading from
	// indexes
	// acquired := zp.doTryLock()

	// if !acquired {
	// 	err = MakeErrConcurrentPageAccess()
	// 	return
	// }

	// defer zp.doUnlock()

	zp.added.Add(*z)
	zp.State = StateChanged

	return
}

func (zp *Page) Flush() (err error) {
	acquired := zp.doTryLock()

	if !acquired {
		err = MakeErrConcurrentPageAccess()
		return
	}

	defer zp.doUnlock()

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

  errors.TodoP3("determine performance of this")
	added := zp.added.Copy()

	if err = added.MergeStream(
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
	acquired := zp.doTryLock()

	if !acquired {
		err = MakeErrConcurrentPageAccess()
		return
	}

	defer zp.doUnlock()

	return zp.copy(w)
}

func (zp *Page) WriteTo(w1 io.Writer) (n int64, err error) {
	acquired := zp.doTryLock()

	if !acquired {
		err = MakeErrConcurrentPageAccess()
		return
	}

	defer zp.doUnlock()

	return zp.writeTo(w1)
}
