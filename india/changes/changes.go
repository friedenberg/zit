package changes

import (
	"strings"

	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/golf/organize_text"
)

type Change struct {
	Etikett, Key string
}

type Changes struct {
	Added   []Change
	Removed []Change
	New     map[string]etikett.Set
}

func ChangesFrom(a1, b1 organize_text.Text) (c Changes) {
	a := a1.ToCompareMap()
	b := b1.ToCompareMap()

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
					Etikett: bez.Etikett,
					Key:     bez.Hinweis,
				},
			)
		}

		delete(a, bez)
	}

	for aez, _ := range a {
		c.Removed = append(
			c.Removed,
			Change{
				Etikett: aez.Etikett,
				Key:     aez.Hinweis,
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
					addNew(z.Bezeichnung, e1)
				}
			}

			// root etiketten
			for _, e2 := range b1.Etiketten() {
				addNew(z.Bezeichnung, e2.String())
			}
		}
	}

	return
}
