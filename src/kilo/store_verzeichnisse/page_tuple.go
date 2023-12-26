package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/kennung_index"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type PageTuple struct {
	All, Schwanzen  Page
	SchwanzenFilter sku.Schwanzen
}

func (pt *PageTuple) initialize(
	n uint8,
	i *Store,
	ki kennung_index.Index,
) {
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
