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
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type Page struct {
	useBestandsaufnahmeForVerzeichnisse bool

	lock *sync.Mutex
	pageId
	schnittstellen.VerzeichnisseFactory
	added       sku.TransactedHeap
	addFilter   schnittstellen.FuncIter[*sku.Transacted]
	flushFilter schnittstellen.FuncIter[*sku.Transacted]

	State
}

func makePage(
	iof schnittstellen.VerzeichnisseFactory, pid pageId,
	fff PageDelegateGetter,
	useBestandsaufnahmeForVerzeichnisse bool,
) (p *Page) {
	flushFilter := collections.MakeWriterNoop[*sku.Transacted]()
	addFilter := collections.MakeWriterNoop[*sku.Transacted]()

	if fff != nil {
		d := fff.GetVerzeichnissePageDelegate(pid.index)

		addFilter = d.ShouldAddVerzeichnisse
		flushFilter = d.ShouldFlushVerzeichnisse
	}

	p = &Page{
		useBestandsaufnahmeForVerzeichnisse: useBestandsaufnahmeForVerzeichnisse,
		lock:                                &sync.Mutex{},
		VerzeichnisseFactory:                iof,
		pageId:                              pid,
		added:                               sku.MakeTransactedHeap(),
		flushFilter:                         flushFilter,
		addFilter:                           addFilter,
	}

	p.added.SetPool(sku.GetTransactedPool())

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

func (zp *Page) Add(z *sku.Transacted) (err error) {
	if z == nil {
		err = errors.Errorf("trying to add nil zettel")
		return
	}

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

	zp.added.Add(z)
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

	// If the cache file does not exist and we have nothing to add, short circuit
	// the flush. This condition occurs on the initial init when the konfig is
	// changed but there are no zettels yet.
	if !files.Exists(zp.path) && zp.added.Len() == 0 {
		return
	}

	if w, err = zp.WriteCloserVerzeichnisse(zp.path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	w1 := bufio.NewWriter(w)

	defer errors.DeferredFlusher(&err, w1)

	m := make(KennungShaMap)

	writeOne := zp.getFuncWriteOne(w1)

	c := iter.MakeChain(
		zp.flushFilter,
		m.ModifyMutter,
		writeOne,
		m.SaveSha,
	)

	if err = zp.copy(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.added.Reset()

	errors.Log().Printf("flushed page: %s", zp.path)

	return
}

func (zp *Page) copy(
	w schnittstellen.FuncIter[*sku.Transacted],
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

	var getOneSku func() (*sku.Transacted, error)

	if zp.useBestandsaufnahmeForVerzeichnisse {
		dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
			r1,
			objekte_format.Default(),
			objekte_format.Options{
				IncludeTai:           true,
				IncludeVerzeichnisse: true,
			},
		)

		getOneSku = func() (sk *sku.Transacted, err error) {
			if !dec.Scan() {
				if err = dec.Error(); err == nil {
					err = io.EOF
				}

				return
			}

			sk = dec.GetTransacted()

			return
		}
	} else {
		dec := gob.NewDecoder(r1)

		getOneSku = func() (sk *sku.Transacted, err error) {
			sk = sku.GetTransactedPool().Get()
			err = dec.Decode(sk)
			return
		}
	}

	errors.TodoP3("determine performance of this")
	added := zp.added.Copy()

	if err = added.MergeStream(
		func() (tz *sku.Transacted, err error) {
			if tz, err = getOneSku(); err != nil {
				if errors.IsEOF(err) {
					err = collections.MakeErrStopIteration()
				} else {
					err = errors.Wrapf(err, "Page: %s", zp.pageId.path)
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

func (zp *Page) getFuncWriteOne(
	w io.Writer,
) schnittstellen.FuncIter[*sku.Transacted] {
	if zp.useBestandsaufnahmeForVerzeichnisse {
		enc := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			w,
			objekte_format.Default(),
			objekte_format.Options{
				IncludeTai:           true,
				IncludeVerzeichnisse: true,
			},
		)

		return func(z *sku.Transacted) (err error) {
			if _, err = enc.Print(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	} else {
		enc := gob.NewEncoder(w)

		return func(z *sku.Transacted) (err error) {
			if err = enc.Encode(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}
}

func (zp *Page) writeTo(w1 io.Writer) (err error) {
	w := bufio.NewWriter(w1)
	defer errors.DeferredFlusher(&err, w)

	writeOne := zp.getFuncWriteOne(w)

	if err = zp.copy(iter.MakeChain(zp.flushFilter, writeOne)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (zp *Page) Copy(
	w schnittstellen.FuncIter[*sku.Transacted],
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
