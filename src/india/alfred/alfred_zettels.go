package alfred

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (w *Writer) zettelToItem(
	z *sku.Transacted,
	ha func(kennung.Hinweis) (string, error),
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.GetMetadatei().Bezeichnung.String()

	es := iter.StringCommaSeparated[kennung.Etikett](
		z.Metadatei.GetEtiketten(),
	)

	if a.Title == "" {
		a.Title = z.GetKennungLike().String()
		a.Subtitle = fmt.Sprintf("%s", es)
	} else {
		a.Subtitle = fmt.Sprintf("%s: %s", z.GetKennungLike().String(), es)
	}

	a.Arg = z.GetKennungLike().String()

	mb := alfred.NewMatchBuilder()

	k := z.GetKennungLike()
	parts := k.Parts()

	mb.AddMatches(k.String())
	mb.AddMatches(parts[0])
	mb.AddMatches(parts[2])
	mb.AddMatches(z.GetMetadatei().Bezeichnung.String())
	mb.AddMatches(z.GetTyp().String())
	z.Metadatei.GetEtiketten().Each(
		func(e kennung.Etikett) (err error) {
			ei, err := w.kennungIndex.GetEtikett(&e)
			if err != nil {
				err = errors.Wrap(err)
				return
			}

			ei.GetExpandedAll().Each(
				func(e kennung.Etikett) (err error) {
					mb.AddMatches(e.String())
					return
				},
			)

			return
		},
	)

	t := z.GetTyp()

	if ti, err := w.typenIndex.Get(&t); err == nil {
		ti.GetExpandedAll().Each(
			func(t kennung.Typ) (err error) {
				mb.AddMatches(t.String())
				return
			},
		)
	}

	// if ha != nil {
	// 	var h hinweis.Hinweis
	// 	var err error

	// 	if h, err = ha.AbbreviateHinweis(z.Hinweis); err != nil {
	// 		return ErrorToItem(err)
	// 	}

	// 	mb.AddMatches(h.String())
	// 	mb.AddMatches(h.Kopf())
	// 	mb.AddMatches(h.Schwanz())
	// }

	a.Match = mb.String()

	// if len(a.Match) > 100 {
	// 	a.Match = a.Match[:100]
	// }

	a.Text.Copy = k.String()
	a.Uid = "zit://" + k.String()

	return
}

func (w *Writer) etikettToItem(
	e *kennung.Etikett,
) (a *alfred.Item) {
	ei, err := w.etikettenIndex.GetEtikett(e)
	a = w.alfredWriter.Get()

	if err != nil {
		a.Subtitle = err.Error()
		return
	} else {
		a.Subtitle = fmt.Sprintf("%d", ei.GetSchwanzenCount())
	}

	a.Title = "@" + e.String()

	a.Arg = e.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(a.Title)
	mb.AddMatches(iter.Strings[kennung.Etikett](kennung.ExpandOne(e))...)

	a.Match = mb.String()

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}

func (w *Writer) errorToItem(err error) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = err.Error()

	return
}

func (w *Writer) hinweisToItem(e kennung.Hinweis) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(),
	// strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatch(e.String())
	mb.AddMatch(e.Kopf())
	mb.AddMatch(e.Schwanz())

	a.Match = mb.String()

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}
