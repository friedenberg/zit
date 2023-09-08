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
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_formats"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type Page struct {
	useBestandsaufnahmeForVerzeichnisse bool

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
	useBestandsaufnahmeForVerzeichnisse bool,
) (p *Page) {
	flushFilter := collections.MakeWriterNoop[*transacted.Zettel]()
	addFilter := collections.MakeWriterNoop[*transacted.Zettel]()

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
		pool:                                pool,
		added:                               zettel.MakeHeapTransacted(),
		flushFilter:                         flushFilter,
		addFilter:                           addFilter,
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

	z.GetMetadateiPtr().Verzeichnisse.ExpandedEtiketten = kennung.ExpandMany[kennung.Etikett](
		z.GetMetadatei().GetEtiketten(),
		kennung.ExpanderRight,
	)

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

	var getOneSku func() (sku.SkuLikePtr, error)

	if zp.useBestandsaufnahmeForVerzeichnisse {
		dec := sku_formats.MakeFormatBestandsaufnahmeScanner(
			r,
			objekte_format.BestandsaufnahmeFormatIncludeTaiVerzeichnisse(),
		)

		getOneSku = func() (sk sku.SkuLikePtr, err error) {
			if !dec.Scan() {
				if err = dec.Error(); err == nil {
					err = io.EOF
				}

				return
			}

			sk = dec.GetSkuLikePtr()

			return
		}
	} else {
		dec := gob.NewDecoder(r)

		getOneSku = func() (sk sku.SkuLikePtr, err error) {
			tz := zp.pool.Get()
			err = dec.Decode(tz)
			sk = tz
			return
		}
	}

	errors.TodoP3("determine performance of this")
	added := zp.added.Copy()

	if err = added.MergeStream(
		func() (tz *transacted.Zettel, err error) {
			var sk sku.SkuLikePtr

			if sk, err = getOneSku(); err != nil {
				if errors.IsEOF(err) {
					err = collections.MakeErrStopIteration()
				} else {
					err = errors.Wrapf(err, "Page: %s", zp.pageId.path)
				}

				return
			}

			ok := false

			if tz, ok = sk.(*transacted.Zettel); !ok {
				err = errors.Errorf("expected %T but got %T, err: %s", tz, sk, err)
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

	var writeOne func(z *transacted.Zettel) error

	if zp.useBestandsaufnahmeForVerzeichnisse {
		enc := sku_formats.MakeFormatBestandsaufnahmePrinter(
			w,
			objekte_format.BestandsaufnahmeFormatIncludeTaiVerzeichnisse(),
		)

		writeOne = func(z *transacted.Zettel) (err error) {
			if _, err = enc.Print(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	} else {
		enc := gob.NewEncoder(w)

		writeOne = func(z *transacted.Zettel) (err error) {
			if err = enc.Encode(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if err = zp.copy(iter.MakeChain(zp.flushFilter, writeOne)); err != nil {
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
