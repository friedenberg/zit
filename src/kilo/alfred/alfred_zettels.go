package alfred

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
)

func (w *Writer) zettelToItem(z *zettel_verzeichnisse.Zettel, ha hinweis.Abbr) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.Transacted.Named.Stored.Objekte.Bezeichnung.String()

	if a.Title == "" {
		a.Title = z.Transacted.Named.Kennung.String()
		a.Subtitle = fmt.Sprintf(
			"%s",
			strings.Join(z.EtikettenSorted, ", "),
		)
	} else {
		a.Subtitle = fmt.Sprintf(
			"%s: %s",
			z.Transacted.Named.Kennung.String(),
			strings.Join(z.EtikettenSorted, ", "),
		)
	}

	a.Arg = z.Transacted.Named.Kennung.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(z.Transacted.Named.Kennung.String())
	mb.AddMatches(z.Transacted.Named.Kennung.Kopf())
	mb.AddMatches(z.Transacted.Named.Kennung.Schwanz())
	mb.AddMatches(z.Transacted.Named.Stored.Objekte.Bezeichnung.String())
	mb.AddMatches(z.Transacted.Named.Stored.Objekte.Typ.String())
	mb.AddMatches(z.EtikettenExpandedSorted...)

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

	a.Text.Copy = z.Transacted.Named.Kennung.String()
	a.Uid = "zit://" + z.Transacted.Named.Kennung.String()

	return
}

func (w *Writer) etikettToItem(e kennung.Etikett) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = "@" + e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(), strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(a.Title)
	mb.AddMatches(e.Expanded().Strings()...)

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

func (w *Writer) hinweisToItem(e hinweis.Hinweis) (a *alfred.Item) {
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
