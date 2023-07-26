package alfred

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

func (w *Writer) zettelToItem(
	z *zettel.Transacted,
	ha func(kennung.Hinweis) (string, error),
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.GetMetadatei().Bezeichnung.String()

	es := iter.StringCommaSeparated[kennung.Etikett](
		z.GetMetadatei().Etiketten,
	)

	if a.Title == "" {
		a.Title = z.Kennung().String()
		a.Subtitle = fmt.Sprintf("%s", es)
	} else {
		a.Subtitle = fmt.Sprintf("%s: %s", z.Kennung().String(), es)
	}

	a.Arg = z.Kennung().String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(z.Kennung().String())
	mb.AddMatches(z.Kennung().Kopf())
	mb.AddMatches(z.Kennung().Schwanz())
	mb.AddMatches(z.GetMetadatei().Bezeichnung.String())
	mb.AddMatches(z.GetTyp().String())
	z.GetMetadatei().Etiketten.Each(
		func(e kennung.Etikett) (err error) {
			ei, err := w.kennungIndex.GetEtikett(e)
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

	if ti, err := w.typenIndex.Get(z.GetTyp()); err == nil {
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

	a.Text.Copy = z.Kennung().String()
	a.Uid = "zit://" + z.Kennung().String()

	return
}

func (w *Writer) etikettToItem(
	ei kennung.IndexedLike[kennung.Etikett, *kennung.Etikett],
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	e := ei.GetKennung()
	a.Title = "@" + e.String()
	a.Subtitle = fmt.Sprintf("%d", ei.GetSchwanzenCount())

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
