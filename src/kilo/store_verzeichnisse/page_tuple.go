package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type PageTuple struct {
	N               uint8
	All, Schwanzen  Page
	SchwanzenFilter sku.Schwanzen
}

func (pt *PageTuple) initialize(
	n uint8,
	i *Store,
	ki kennung_index.Index,
) {
	pt.N = n

	pt.All.initialize(
		i.standort.SansAge().SansCompression(),
		i.PageIdForIndex(n, false),
		i.ennui,
	)

	pt.SchwanzenFilter.Initialize(ki, i.applyKonfig)

	pt.Schwanzen.initializeWithSchwanzen(
		i.standort.SansAge().SansCompression(),
		i.PageIdForIndex(uint8(n), true),
		&pt.SchwanzenFilter,
	)
}

func (pt *PageTuple) Add(
	z1 *sku.Transacted,
) (err error) {
	z := sku.GetTransactedPool().Get()

	if err = z.SetFromSkuLike(z1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sku.CalculateAndSetSha(z, nil, objekte_format.Options{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pt.All.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pt.Schwanzen.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (pt *PageTuple) PageForSigil(s kennung.Sigil) *Page {
	if s.IncludesHistory() {
		return &pt.All
	} else {
		return &pt.Schwanzen
	}
}

func (pt *PageTuple) SetNeedsFlush() {
	pt.All.State = StateChanged
	pt.Schwanzen.State = StateChanged
}

func (pt *PageTuple) Flush() (err error) {
	m := make(KennungShaMap)

	if err = pt.All.Flush(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = pt.Schwanzen.Flush(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type ShaTuple struct {
	Sha, Mutter *sha.Sha
}

type KennungShaMap map[string]ShaTuple

func (ksm KennungShaMap) ReadMutter(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	old := ksm[k.String()]

	if old.Mutter.IsNull() {
		return
	}

	if err = z.GetMetadatei().Mutter.SetShaLike(old.Mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ksm KennungShaMap) SaveSha(z *sku.Transacted) (err error) {
	k := z.GetKennung()
	var sh sha.Sha

	if err = sh.SetShaLike(&z.GetMetadatei().Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	old := ksm[k.String()]
	old.Mutter = old.Sha
	old.Sha = &sh
	ksm[k.String()] = old

	return
}
