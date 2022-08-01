package changes

import (
	"github.com/friedenberg/zit/alfa/logz"
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

	logz.Printf("%#v", a)
	logz.Printf("%#v", b)

	c.Added = make([]Change, 0)
	c.Removed = make([]Change, 0)

	for tuple, _ := range b.Named {
		if _, ok := a.Named[tuple]; ok {
			//zettel had etikett previously
		} else {
			//zettel did not have etikett previously
			c.Added = append(
				c.Added,
				Change{
					Etikett: tuple.Etikett,
					Key:     tuple.Key,
				},
			)
		}

		delete(a.Named, tuple)
	}

	for tuple, _ := range a.Named {
		c.Removed = append(
			c.Removed,
			Change{
				Etikett: tuple.Etikett,
				Key:     tuple.Key,
			},
		)
	}

	c.New = make(map[string]etikett.Set)

	addNew := func(bez, ett string) {
		existing, ok := c.New[bez]

		if !ok {
			existing = etikett.MakeSet()
		}

		existing.AddString(ett)
		c.New[bez] = existing
	}

	for tuple, _ := range b.Unnamed {
		addNew(tuple.Key, tuple.Etikett)
	}

	return
}
