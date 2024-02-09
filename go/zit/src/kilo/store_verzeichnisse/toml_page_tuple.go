package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/objekte_mode"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

type TomlPageTuple struct {
	PageId
	// All, Schwanzen  Page
	ennuiShas, ennuiKennung ennui.Ennui
	added, addedSchwanz     *sku.TransactedHeap
	flushMode               objekte_mode.Mode
	hasChanges              bool
	changesAreHistorical    bool
	standort                standort.Standort
	konfig                  *konfig.Compiled
	etikettIndex            kennung_index.EtikettIndexMutation
}

func (pt *TomlPageTuple) initialize(
	pid PageId,
	i *Store,
	ki kennung_index.Index,
) {
	pt.standort = i.standort.SansAge().SansCompression()
	pt.PageId = pid
	pt.added = sku.MakeTransactedHeap()
	pt.addedSchwanz = sku.MakeTransactedHeap()
	pt.etikettIndex = ki
	pt.ennuiKennung = i.ennuiKennung
	pt.konfig = i.erworben
}

func (pt *TomlPageTuple) add(
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

func (pt *TomlPageTuple) waitingToAddLen() int {
	return pt.added.Len() + pt.addedSchwanz.Len()
}

func (pt *TomlPageTuple) SetNeedsFlushHistory() {
	pt.hasChanges = true
	pt.changesAreHistorical = true
}

func (pt *TomlPageTuple) CopyEverything(
	s kennung.Sigil,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeSchwanz(s, w, true, true)
}

func (pt *TomlPageTuple) CopyJustHistory(
	s kennung.Sigil,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeSchwanz(s, w, false, false)
}

func (pt *TomlPageTuple) CopyJustHistoryFrom(
	r io.Reader,
	s kennung.Sigil,
	w schnittstellen.FuncIter[sku_fmt.Sku],
) (err error) {
	dec := sku_fmt.Binary{Sigil: s}

	var sk sku_fmt.Sku

	for {
		sk.Offset += sk.ContentLength
		sk.Transacted = sku.GetTransactedPool().Get()
		sk.ContentLength, err = dec.ReadFormatAndMatchSigil(r, &sk)

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

func (pt *TomlPageTuple) CopyJustHistoryAndAdded(
	s kennung.Sigil,
	w schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return pt.copyHistoryAndMaybeSchwanz(s, w, true, false)
}

func (pt *TomlPageTuple) copyHistoryAndMaybeSchwanz(
	s kennung.Sigil,
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
			s,
			func(sk sku_fmt.Sku) (err error) {
				return w(sk.Transacted)
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	dec := sku_fmt.Binary{Sigil: s}

	errors.TodoP3("determine performance of this")
	added := pt.added.Copy()

	var sk sku_fmt.Sku

	if err = added.MergeStream(
		func() (tz *sku.Transacted, err error) {
			tz = sku.GetTransactedPool().Get()
			sk.Transacted = tz
			_, err = dec.ReadFormatAndMatchSigil(br, &sk)

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

	if err = addedSchwanz.MergeStream(
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

func (pt *TomlPageTuple) Flush() (err error) {
	pw := &tomlPageWriter{
		TomlPageTuple: pt,
	}

	if err = pw.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	pt.hasChanges = false
	pt.changesAreHistorical = false

	return
}
