package store_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

func (i *Zettelen) addZettelHinweis(tz zettel_transacted.Zettel) (err error) {
	var n int

	if n, err = i.PageForHinweis(tz.Named.Hinweis); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	z := i.pool.Get()
	z.PageSelection.Reason = PageSelectionReasonHinweis
	z.Transacted = tz
	z.EtikettenExpandedSorted = tz.Named.Stored.Zettel.Etiketten.Expanded().SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) addZettelTransacted(tz zettel_transacted.Zettel) (err error) {
	var n int

	if n, err = i.PageForTransacted(tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	z := i.pool.Get()
	z.PageSelection.Reason = PageSelectionReasonStoredSha
	z.Transacted = tz
	z.EtikettenExpandedSorted = tz.Named.Stored.Zettel.Etiketten.Expanded().SortedString()
	z.EtikettenSorted = tz.Named.Stored.Zettel.Etiketten.SortedString()

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
