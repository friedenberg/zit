package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/delta/heap"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/konfig"
)

type PageId = sha.PageId

type Page struct {
	PageId
	ennuiStore
	// All, Schwanzen  Page
	added, addedSchwanz *sku.TransactedHeap
	flushMode           objekte_mode.Mode
	hasChanges          bool
	standort            standort.Standort
	konfig              *konfig.Compiled
}

func (pt *Page) initialize(
	pid PageId,
	i *Store,
) {
	pt.standort = i.standort.SansAge().SansCompression()
	pt.PageId = pid
	pt.added = sku.MakeTransactedHeap()
	pt.addedSchwanz = sku.MakeTransactedHeap()
	pt.konfig = i.erworben
	pt.ennuiStore = i.ennuiStore
}

func (pt *Page) add(
	z1 *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	z := sku.GetTransactedPool().Get()

	if err = z.SetFromSkuLike(z1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mode.Contains(objekte_mode.ModeSchwanz) {
		pt.addedSchwanz.Add(z)
	} else {
		pt.added.Add(z)
	}

	pt.hasChanges = true

	return
}

func (pt *Page) waitingToAddLen() int {
	return pt.added.Len() + pt.addedSchwanz.Len()
}

func (pt *Page) CopyJustHistory(
	s sku.QueryGroup,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeSchwanz(s, w, false, false)
}

func (pt *Page) CopyJustHistoryFrom(
	r io.Reader,
	s sku.QueryGroup,
	w schnittstellen.FuncIter[skuWithRangeAndSigil],
) (err error) {
	dec := makeBinaryWithQueryGroup(s, kennung.SigilHistory)

	var sk skuWithRangeAndSigil

	for {
		sk.Offset += sk.ContentLength
		sk.Transacted = sku.GetTransactedPool().Get()
		sk.ContentLength, err = dec.readFormatAndMatchSigil(r, &sk)
		if err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if err = w(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}
}

func (pt *Page) CopyJustHistoryAndAdded(
	s sku.QueryGroup,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeSchwanz(s, w, true, false)
}

func (pt *Page) copyHistoryAndMaybeSchwanz(
	qg sku.QueryGroup,
	w schnittstellen.FuncIter[*sku.Transacted],
	includeAdded bool,
	includeAddedSchwanz bool,
) (err error) {
	var r io.ReadCloser

	if r, err = pt.standort.ReadCloserVerzeichnisse(pt.Path()); err != nil {
		if errors.IsNotExist(err) {
			r = io.NopCloser(bytes.NewReader(nil))
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, r)

	br := bufio.NewReader(r)

	if !includeAdded && !includeAddedSchwanz {
		if err = pt.CopyJustHistoryFrom(
			br,
			qg,
			func(sk skuWithRangeAndSigil) (err error) {
				if err = w(sk.Transacted); err != nil {
					err = errors.Wrapf(err, "%s", sk.Transacted)
					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	dec := makeBinaryWithQueryGroup(qg, kennung.SigilHistory)

	errors.TodoP3("determine performance of this")
	added := pt.added.Copy()

	var sk skuWithRangeAndSigil

	if err = heap.MergeStream(
		&added,
		func() (tz *sku.Transacted, err error) {
			tz = sku.GetTransactedPool().Get()
			sk.Transacted = tz
			_, err = dec.readFormatAndMatchSigil(br, &sk)
			if err != nil {
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

	if !includeAddedSchwanz {
		return
	}

	addedSchwanz := pt.addedSchwanz.Copy()

	if err = heap.MergeStream(
		&addedSchwanz,
		func() (tz *sku.Transacted, err error) {
			err = collections.MakeErrStopIteration()
			return
		},
		w,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (pt *Page) MakeFlush(changesAreHistorical bool) func() error {
	return func() (err error) {
		pw := &writer{
			Page: pt,
		}

		if changesAreHistorical {
			pw.changesAreHistorical = true
			pw.hasChanges = true
		}

		if err = pw.Flush(); err != nil {
			err = errors.Wrap(err)
			return
		}

		pt.hasChanges = false

		return
	}
}
