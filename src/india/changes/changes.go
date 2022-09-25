package changes

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/hotel/organize_text"
)

type Change struct {
	Etikett, Key string
}

type New struct {
	Key       string
	Etiketten etikett.Set
}

type Changes struct {
	Added   []Change
	Removed []Change
	New     []New
	AllB    []string
}

type changes struct {
	Added   []Change
	Removed []Change
	New     map[string]etikett.Set
	AllB    []string
}

func ChangesFrom(a1, b1 *organize_text.Text) (c1 Changes, err error) {
	var c changes
	var a, b organize_text.CompareMap

	if a, err = a1.ToCompareMap(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if b, err = b1.ToCompareMap(); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Added = make([]Change, 0)
	c.Removed = make([]Change, 0)
	c.AllB = make([]string, 0, len(b.Named))

	for h, es1 := range b.Named {
		c.AllB = append(c.AllB, h)

		for _, e1 := range es1 {
			if a.Named.Contains(h, e1) {
				//zettel had etikett previously
			} else {
				c.Added = append(
					c.Added,
					Change{
						Etikett: e1.String(),
						Key:     h,
					},
				)
			}

			if es2, ok := a.Named[h]; ok {
				es2.Remove(e1)
				a.Named[h] = es2
			}
		}
	}

	for h, es := range a.Named {
		for _, e1 := range es {
			if e1.String() == "" {
				err = errors.Errorf("empty etikett for %s", h)
				return
			}

			c.Removed = append(
				c.Removed,
				Change{
					Etikett: e1.String(),
					Key:     h,
				},
			)
		}
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

	for h, es := range b.Unnamed {
		for _, e := range es {
			addNew(h, e.String())
		}
	}

	c1 = c.toChanges()

	return
}

func (c changes) toChanges() (c1 Changes) {
	c1.Added = c.Added
	c1.Removed = c.Removed
	c1.AllB = c.AllB
	c1.New = make([]New, 0, len(c1.New))

	for h, e := range c.New {
		c1.New = append(c1.New, New{Key: h, Etiketten: e})
	}

	sort.Slice(
		c1.Added,
		func(i, j int) bool { return c1.Added[i].Key < c1.Added[j].Key },
	)

	sort.Slice(
		c1.Removed,
		func(i, j int) bool { return c1.Removed[i].Key < c1.Removed[j].Key },
	)

	sort.Slice(
		c1.AllB,
		func(i, j int) bool { return c1.AllB[i] < c1.AllB[j] },
	)

	sort.Slice(
		c1.New,
		func(i, j int) bool { return c1.New[i].Key < c1.New[j].Key },
	)

	return
}
