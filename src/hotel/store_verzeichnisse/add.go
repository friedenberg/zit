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

	z := i.MakeZettel(
		tz,
		PageSelectionReasonHinweis,
		"",
	)

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

	z := i.MakeZettel(
		tz,
		PageSelectionReasonStoredSha,
		"",
	)

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) addZettelAkte(tz zettel_transacted.Zettel) (err error) {
	var n int

	if n, err = i.PageForSha(tz.Named.Stored.Zettel.Akte); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.ValidatePageIndex(n); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := i.pages[n]

	z := i.MakeZettel(
		tz,
		PageSelectionReasonAkte,
		"",
	)

	if err = p.Add(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Zettelen) addZettelEtikett(tz zettel_transacted.Zettel) (err error) {
	z := i.MakeZettel(
		tz,
		PageSelectionReasonEtikett,
		"",
	)

	for _, e := range tz.Named.Stored.Zettel.Etiketten.Etiketten() {
		var n int

		if n, err = i.PageForEtikett(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = i.ValidatePageIndex(n); err != nil {
			err = errors.Wrap(err)
			return
		}

		p := i.pages[n]

		if err = p.Add(z); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
