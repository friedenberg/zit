package alfred

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/alfred"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

func ZettelToItem(z zettel_named.Zettel, ha hinweis.Abbr) (a alfred.Item) {
	a.Title = z.Stored.Zettel.Bezeichnung.String()

	if a.Title == "" {
		a.Title = z.Hinweis.String()
		a.Subtitle = fmt.Sprintf(
			"%s",
			strings.Join(z.Stored.Zettel.Etiketten.SortedString(), ", "),
		)
	} else {
		a.Subtitle = fmt.Sprintf(
			"%s: %s",
			z.Hinweis.String(),
			strings.Join(z.Stored.Zettel.Etiketten.SortedString(), ", "),
		)
	}

	a.Arg = z.Hinweis.String()

	mb := alfred.NewMatchBuilder()

	mb.AddMatches(z.Hinweis.String())
	mb.AddMatches(z.Hinweis.Kopf())
	mb.AddMatches(z.Hinweis.Schwanz())
	mb.AddMatches(z.Stored.Zettel.Bezeichnung.String())
	mb.AddMatches(z.Stored.Zettel.Typ.String())
	mb.AddMatches(EtikettenStringsFromZettel(z.Stored.Zettel.Etiketten, true)...)

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

	a.Text.Copy = z.Hinweis.String()
	a.Uid = "zit://" + z.Hinweis.String()

	return
}

func EtikettToItem(e etikett.Etikett) (a alfred.Item) {
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

func ErrorToItem(err error) (a alfred.Item) {
	a.Title = err.Error()

	return
}

func HinweisToItem(e hinweis.Hinweis) (a alfred.Item) {
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

func EtikettenStringsFromZettel(es etikett.Set, shouldExpand bool) (out []string) {
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
