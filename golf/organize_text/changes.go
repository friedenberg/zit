package organize_text

import (
	"strings"
)

type Change struct {
	Etikett, Hinweis string
}

type Changes struct {
	Added   []Change
	Removed []Change
}

func (a1 *organizeText) ChangesFrom(b1 Text) (c Changes) {
	type etikettZettel struct {
		etikett, hinweis string
	}

	type compareMap map[etikettZettel]bool

	makeCompareMap := func(in Text) (out compareMap) {
		out = make(compareMap)

		for e, zs := range in.Zettels() {
			for z, _ := range zs {
				// individual etiketten
				for _, e1 := range strings.Split(e, ", ") {
					// root etiketten have an empty string representation
					if e1 != "" {
						out[etikettZettel{etikett: e1, hinweis: z.hinweis}] = true
					}
				}

				// root etiketten
				for _, e2 := range in.Etiketten() {
					out[etikettZettel{etikett: e2.String(), hinweis: z.hinweis}] = true
				}
			}
		}

		return
	}

	a := makeCompareMap(a1)
	b := makeCompareMap(b1)

	c.Added = make([]Change, 0)
	c.Removed = make([]Change, 0)

	for bez, _ := range b {
		if _, ok := a[bez]; ok {
			//zettel had etikett previously
		} else {
			//zettel did not have etikett previously
			c.Added = append(
				c.Added,
				Change{
					Etikett: bez.etikett,
					Hinweis: bez.hinweis,
				},
			)
		}

		delete(a, bez)
	}

	for aez, _ := range a {
		c.Removed = append(
			c.Removed,
			Change{
				Etikett: aez.etikett,
				Hinweis: aez.hinweis,
			},
		)
	}

	return
}
