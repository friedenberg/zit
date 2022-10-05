package alfred

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

func (w *Writer) zettelToItem(z zettel_transacted.Zettel, ha hinweis.Abbr) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = z.Named.Stored.Zettel.Bezeichnung.String()

	if a.Title == "" {
		a.Title = z.Named.Hinweis.String()
		a.Subtitle = fmt.Sprintf(
			"%s",
			strings.Join(z.Named.Stored.Zettel.Etiketten.SortedString(), ", "),
		)
	} else {
		a.Subtitle = fmt.Sprintf(
			"%s: %s",
			z.Named.Hinweis.String(),
			strings.Join(z.Named.Stored.Zettel.Etiketten.SortedString(), ", "),
		)
	}

	a.Arg = z.Named.Hinweis.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(z.Named.Hinweis.String())
	mb.AddMatches(z.Named.Hinweis.Kopf())
	mb.AddMatches(z.Named.Hinweis.Schwanz())
	mb.AddMatches(z.Named.Stored.Zettel.Bezeichnung.String())
	mb.AddMatches(z.Named.Stored.Zettel.Typ.String())
	mb.AddMatches(w.etikettenStringsFromZettel(z.Named.Stored.Zettel.Etiketten, true)...)

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

	a.Text.Copy = z.Named.Hinweis.String()
	a.Uid = "zit://" + z.Named.Hinweis.String()

	return
}

func (w *Writer) etikettToItem(e etikett.Etikett) (a *alfred.Item) {
	a = w.alfredWriter.Get()

	a.Title = e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(), strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := alfred.NewMatchBuilder()

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

func (w *Writer) etikettenStringsFromZettel(es etikett.Set, shouldExpand bool) (out []string) {
	out = make([]string, 0, es.Len())

	for _, e := range es.Etiketten() {
		if shouldExpand {
			out = append(out, e.Expanded().Strings()...)
		} else {
			out = append(out, e.String())
		}
	}

	return
}
