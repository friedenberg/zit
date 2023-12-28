package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

type Page struct {
	pageId
	schnittstellen.VerzeichnisseFactory
	added       *sku.TransactedHeap
	addFilter   schnittstellen.FuncIter[*sku.Transacted]
	flushFilter schnittstellen.FuncIter[*sku.Transacted]
	delegate    PageDelegate
	ennui       ennui.Ennui

	State
}

func (p *Page) initialize(
	iof schnittstellen.VerzeichnisseFactory,
	pid pageId,
	e ennui.Ennui,
) {
	p.VerzeichnisseFactory = iof
	p.pageId = pid
	p.added = sku.MakeTransactedHeap()
	p.addFilter = collections.MakeWriterNoop[*sku.Transacted]()
	p.flushFilter = collections.MakeWriterNoop[*sku.Transacted]()
	p.ennui = e

	p.added.SetPool(sku.GetTransactedPool())
}

func (p *Page) initializeWithSchwanzen(
	iof schnittstellen.VerzeichnisseFactory,
	pid pageId,
	d PageDelegate,
) {
	p.initialize(iof, pid, nil)

	if d == nil {
		return
	}

	p.delegate = d

	p.addFilter = p.delegate.ShouldAddVerzeichnisse
	p.flushFilter = p.delegate.ShouldFlushVerzeichnisse
}

func (zp *Page) Add(z *sku.Transacted) (err error) {
	if z == nil {
		panic("trying to add nil zettel")
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

	zp.added.Add(z)
	zp.State = StateChanged

	log.Log().Printf("added %s", z.Kennung.String())

	return
}

func (zp *Page) Flush(m KennungShaMap) (err error) {
	if zp.State < StateChanged {
		errors.Log().Printf("not flushing page: %s", zp.path)
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

	writeOne := zp.getFuncWriteOne(w1)

	c := iter.MakeChain(
		zp.flushFilter,
		m.ReadMutter,
		writeOne,
		m.SaveSha,
	)

	if err = zp.copy(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	zp.added.Reset()

	// errors.Log().Printf("flushed page: %s", zp.path)

	return
}

func (zp *Page) copy(
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	var r1 io.ReadCloser

	if r1, err = zp.ReadCloserVerzeichnisse(zp.path); err != nil {
		if errors.IsNotExist(err) {
			r1 = io.NopCloser(bytes.NewReader(nil))
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var getOneSku func() (*sku.Transacted, error)

	dec := sku_fmt.MakeFormatBestandsaufnahmeScanner(
		r1,
		objekte_format.Default(),
		options,
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

	errors.TodoP3("determine performance of this")
	added := zp.added.Copy()

	if err = added.MergeStream(
		func() (tz *sku.Transacted, err error) {
			if tz, err = getOneSku(); err != nil {
				if errors.IsEOF(err) {
					err = collections.MakeErrStopIteration()
				} else {
					err = errors.Wrapf(err, "Page: %s", zp.path)
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
	enc := sku_fmt.MakeFormatBestandsaufnahmePrinter(
		w,
		objekte_format.Default(),
		options,
	)

	return func(z *sku.Transacted) (err error) {
		offset := enc.Offset()

		if _, err = enc.Print(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if zp.ennui == nil {
			return
		}

		if err = zp.ennui.AddMetadatei(
			z.GetMetadatei(),
			zp.index,
			uint64(offset),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
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
	return zp.copy(w)
}

func (zp *Page) WriteTo(w1 io.Writer) (n int64, err error) {
	err = zp.writeTo(w1)

	return
}
