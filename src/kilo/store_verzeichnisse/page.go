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
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/transacted"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type Page struct {
	lock *sync.Mutex
	pageId
	schnittstellen.VerzeichnisseFactory
	pool        schnittstellen.Pool[transacted.Zettel, *transacted.Zettel]
	added       zettel.HeapTransacted
	addFilter   schnittstellen.FuncIter[*transacted.Zettel]
	flushFilter schnittstellen.FuncIter[*transacted.Zettel]
	State
}

func makeZettelenPage(
	iof schnittstellen.VerzeichnisseFactory, pid pageId,
	pool schnittstellen.Pool[transacted.Zettel, *transacted.Zettel],
	fff PageDelegateGetter,
) (p *Page) {
	flushFilter := collections.MakeWriterNoop[*transacted.Zettel]()
	addFilter := collections.MakeWriterNoop[*transacted.Zettel]()

	if fff != nil {
		d := fff.GetVerzeichnissePageDelegate(pid.index)

		addFilter = d.ShouldAddVerzeichnisse
		flushFilter = d.ShouldFlushVerzeichnisse
	}

	p = &Page{
		lock:                 &sync.Mutex{},
		VerzeichnisseFactory: iof,
		pageId:               pid,
		pool:                 pool,
		added:                zettel.MakeHeapTransacted(),
		flushFilter:          flushFilter,
		addFilter:            addFilter,
	}

	p.added.SetPool(p.pool)

	return
}

func (zp *Page) doTryLock() (ok bool) {
	ok = zp.lock.TryLock()

	if ok {
	} else {
	}

	return
}

func (zp *Page) doLock() {
	zp.lock.Lock()
}

func (zp *Page) doUnlock() {
	zp.lock.Unlock()
}

func (zp *Page) Add(z *transacted.Zettel) (err error) {
	if z == nil {
		err = errors.Errorf("trying to add nil zettel")
		return
	}

	log.Log().Printf("adding: %s", z.GetSkuLike())

	if err = zp.addFilter(z); err != nil {
		if iter.IsStopIteration(err) {
			errors.Log().Printf("eliding %s", z.Kennung)
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

	log.Log().Printf("added: %s", z.GetSkuLike())

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

	if err = zp.writeTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.added.Reset()

	errors.Log().Printf("flushed page: %s", zp.path)

	return
}

func (zp *Page) copy(
	w schnittstellen.FuncIter[*transacted.Zettel],
) (err error) {
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
		func() (tz *transacted.Zettel, err error) {
			tz = zp.pool.Get()

			if err = dec.Decode(tz); err != nil {
				if errors.IsEOF(err) {
					err = collections.MakeErrStopIteration()
				} else {
					err = errors.Wrap(err)
				}

				return
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

func (zp *Page) writeTo(w1 io.Writer) (err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	enc := gob.NewEncoder(w)

	if err = zp.copy(
		iter.MakeChain(
			zp.flushFilter,
			func(z *transacted.Zettel) (err error) {
				return enc.Encode(z)
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *Page) Copy(
	w schnittstellen.FuncIter[*transacted.Zettel],
) (err error) {
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

	err = zp.writeTo(w1)

	return
}
