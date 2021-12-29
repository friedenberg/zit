package alfred_zettels

import (
	"fmt"
	"strings"
)

func ZettelToItem(z _NamedZettel) (a _AlfredItem) {
	a.Title = z.Zettel.Bezeichnung.String()

	if a.Title == "" {
		a.Title = z.Hinweis.String()
		a.Subtitle = fmt.Sprintf(
			"%s",
			strings.Join(z.Zettel.Etiketten.SortedString(), ", "),
		)
	} else {
		a.Subtitle = fmt.Sprintf(
			"%s: %s",
			z.Hinweis.String(),
			strings.Join(z.Zettel.Etiketten.SortedString(), ", "),
		)
	}

	a.Arg = z.Hinweis.String()

	mb := _AlfredNewMatchBuilder()

	mb.AddMatches(z.Hinweis.String())
	mb.AddMatches(z.Hinweis.Head())
	mb.AddMatches(z.Hinweis.Tail())
	mb.AddMatches(z.Zettel.Bezeichnung.String())
	mb.AddMatches(EtikettenStringsFromZettel(z.Zettel.Etiketten, true)...)

	a.Match = mb.String()

	// if len(a.Match) > 100 {
	// 	a.Match = a.Match[:100]
	// }

	a.Text.Copy = z.Hinweis.String()
	a.Uid = "zit://" + z.Hinweis.String()

	return
}

func EtikettToItem(e _Etikett) (a _AlfredItem) {
	a.Title = e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(), strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := _AlfredNewMatchBuilder()

	mb.AddMatches(e.Expanded().Strings()...)

	a.Match = mb.String()

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}

func HinweisToItem(e _Hinweis) (a _AlfredItem) {
	a.Title = e.String()
	// a.Subtitle = fmt.Sprintf("%s: %s", z.Hinweis.String(), strings.Join(EtikettenStringsFromZettel(z, false), ", "))

	a.Arg = e.String()

	mb := _AlfredNewMatchBuilder()

	mb.AddMatch(e.String())
	mb.AddMatch(e.Head())
	mb.AddMatch(e.Tail())

	a.Match = mb.String()

	a.Text.Copy = e.String()
	a.Uid = "zit://" + e.String()

	return
}

func EtikettenStringsFromZettel(es _EtikettSet, shouldExpand bool) (out []string) {
	out = make([]string, 0, es.Len())

	for _, e := range es {
		if shouldExpand {
			out = append(out, e.Expanded().Strings()...)
		} else {
			out = append(out, e.String())
		}
	}

	return
}
