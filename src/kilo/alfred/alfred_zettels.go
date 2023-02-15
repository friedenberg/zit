package alfred

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

func (w *Writer) zettelToItem(
	z *zettel.Transacted,
	ha schnittstellen.FuncAbbreviateKorper,
) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.Objekte.Bezeichnung.String()

	if a.Title == "" {
		a.Title = z.Kennung().String()
		a.Subtitle = fmt.Sprintf(
			"%s",
			strings.Join(z.Verzeichnisse.Etiketten.Sorted, ", "),
		)
	} else {
		a.Subtitle = fmt.Sprintf(
			"%s: %s",
			z.Kennung().String(),
			strings.Join(z.Verzeichnisse.Etiketten.Sorted, ", "),
		)
	}

	a.Arg = z.Kennung().String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(z.Kennung().String())
	mb.AddMatches(z.Kennung().Kopf())
	mb.AddMatches(z.Kennung().Schwanz())
	mb.AddMatches(z.Objekte.Bezeichnung.String())
	mb.AddMatches(z.Objekte.Typ.String())
	mb.AddMatches(z.Verzeichnisse.Etiketten.SortedExpanded...)
	mb.AddMatches(z.Verzeichnisse.Typ.Expanded...)

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

func (w *Writer) etikettToItem(e kennung.Etikett) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = "@" + e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(), strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(a.Title)
	mb.AddMatches(collections.Strings(e.Expanded())...)

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
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(), strings.Join(EtikettenStringsFromZettel(z, false), ", "))

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
