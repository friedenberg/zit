package organize_text

import (
	"strings"

	"github.com/friedenberg/zit/charlie/etikett"
)

type Change struct {
	Etikett, Key string
}

type Changes struct {
	Added   []Change
	Removed []Change
	New     map[string]etikett.Set
}

func (a1 *organizeText) ChangesFrom(b1 Text) (c Changes) {
	type etikettZettel struct {
		etikett, hinweis string
	}

	type compareMap map[etikettZettel]bool

	makeCompareMap := func(in Text) (out compareMap) {
		out = make(compareMap)

		for e, zs := range in.ZettelsExisting() {
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
					Key:     bez.hinweis,
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
				Key:     aez.hinweis,
			},
		)
	}

	c.New = make(map[string]etikett.Set)

	addNew := func(bez, ett string) {
		existing, ok := c.New[bez]

		if !ok {
			existing = etikett.NewSet()
		}

		existing.AddString(ett)
		c.New[bez] = existing
	}

	for e, zs := range b1.ZettelsNew() {
		for z, _ := range zs {
			// individual etiketten
			for _, e1 := range strings.Split(e, ", ") {
				// root etiketten have an empty string representation
				if e1 != "" {
					addNew(z.bezeichnung, e1)
				}
			}

			// root etiketten
			for _, e2 := range b1.Etiketten() {
				addNew(z.bezeichnung, e2.String())
			}
		}
	}

	return
}
