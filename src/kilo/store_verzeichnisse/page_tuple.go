package store_verzeichnisse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"

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

type PageId struct {
	Index uint8
	Dir   string
}

func (pid *PageId) PathForSigil(s kennung.Sigil) string {
	if s.IncludesHistory() {
		return pid.PathForAll()
	} else {
		return pid.PathForSchwanz()
	}
}

func (pid *PageId) PathForAll() string {
	return filepath.Join(pid.Dir, fmt.Sprintf("All-%x", pid.Index))
}

func (pid *PageId) PathForSchwanz() string {
	return filepath.Join(pid.Dir, fmt.Sprintf("Schwanz-%x", pid.Index))
}

func (pid *PageId) Path() string {
	return filepath.Join(pid.Dir, fmt.Sprintf("%x", pid.Index))
}

type PageTuple struct {
	PageId
	// All, Schwanzen  Page
	ennuiShas, ennuiKennung ennui.Ennui
	added, addedSchwanz     *sku.TransactedHeap
	hasChanges              bool
	standort                standort.Standort
	konfig                  *konfig.Compiled
	etikettIndex            kennung_index.EtikettIndexMutation
}

func (pt *PageTuple) initialize(
	pid PageId,
	i *Store,
	ki kennung_index.Index,
) {
	pt.standort = i.standort.SansAge().SansCompression()
	pt.PageId = pid
	pt.added = sku.MakeTransactedHeap()
	pt.addedSchwanz = sku.MakeTransactedHeap()
	pt.etikettIndex = ki
	pt.ennuiShas = i.ennuiShas
	pt.ennuiKennung = i.ennuiKennung
	pt.konfig = i.erworben
}

func (pt *PageTuple) add(
	z1 *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	z := sku.GetTransactedPool().Get()

	if err = z.SetFromSkuLike(z1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mode.Contains(objekte_mode.ModeSchwanz) {
		// pt.addedSchwanz.Add(z)
		pt.added.Add(z)
	} else {
		pt.added.Add(z)
	}

	pt.hasChanges = true

	return
}

func (pt *PageTuple) waitingToAddLen() int {
	return pt.added.Len() + pt.addedSchwanz.Len()
}

func (pt *PageTuple) SetNeedsFlush() {
	pt.hasChanges = true
}

func (pt *PageTuple) Copy(
	s kennung.Sigil,
	w schnittstellen.FuncIter[*sku.Transacted],
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

	br := bufio.NewReader(r)

	dec := sku_fmt.Binary{Sigil: s}

	errors.TodoP3("determine performance of this")
	added := pt.added.Copy()

	if err = added.MergeStream(
		func() (tz *sku.Transacted, err error) {
			tz = sku.GetTransactedPool().Get()
			_, err = dec.ReadFormatAndMatchSigil(br, tz)

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

func (pt *PageTuple) Flush() error {
	pw := pageWriter{
		PageTuple: pt,
	}

	return pw.flush()
}
