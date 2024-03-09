package changes2

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Changeable interface {
	CompareMap(
		hinweis_expander func(string) (*kennung.Hinweis, error),
	) (out CompareMap, err error)
	GetSkus(original sku.TransactedSet) (sku.TransactedSet, error)
}

type Changes interface {
	GetChanges() (self Changes, a, b Changeable)
	GetCompareMaps() (a, b CompareMap)
	GetModified() schnittstellen.SetLike[*ChangeBezeichnung]
	GetExisting() schnittstellen.SetLike[*Change]
	GetAddedUnnamed() schnittstellen.SetLike[*Change]
	GetAddedNamed() schnittstellen.SetLike[*Change]
	GetAllBKeys() schnittstellen.SetLike[values.String]
}

type changes struct {
	a, b                     Changeable
	compareMapA, compareMapB CompareMap
	modified                 schnittstellen.MutableSetLike[*ChangeBezeichnung]
	existing                 schnittstellen.MutableSetLike[*Change]
	addedUnnamed             schnittstellen.MutableSetLike[*Change]
	addedNamed               schnittstellen.MutableSetLike[*Change]
	allB                     schnittstellen.MutableSetLike[values.String]
}

func (c changes) GetChanges() (Changes, Changeable, Changeable) {
	return c, c.a, c.b
}

func (c changes) GetCompareMaps() (a, b CompareMap) {
	return c.compareMapA, c.compareMapB
}

func (c changes) GetModified() schnittstellen.SetLike[*ChangeBezeichnung] {
	return c.modified
}

func (c changes) GetExisting() schnittstellen.SetLike[*Change] {
	return c.existing
}

func (c changes) GetAddedUnnamed() schnittstellen.SetLike[*Change] {
	return c.addedUnnamed
}

func (c changes) GetAddedNamed() schnittstellen.SetLike[*Change] {
	return c.addedNamed
}

func (c changes) GetAllBKeys() schnittstellen.SetLike[values.String] {
	return c.allB
}

func ChangesFrom(
	a1, b1 Changeable,
	hinweis_expander func(string) (*kennung.Hinweis, error),
) (c1 Changes, err error) {
	var c changes
	c1 = &c

	c.a, c.b = a1, b1

	if c.compareMapA, err = a1.CompareMap(hinweis_expander); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.compareMapB, err = b1.CompareMap(hinweis_expander); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.modified = collections_value.MakeMutableSet[*ChangeBezeichnung](
		ChangeBezeichnungKeyer{},
	)
	c.existing = collections_value.MakeMutableValueSet[*Change](
		ChangeKeyer{},
	)
	c.addedUnnamed = collections_value.MakeMutableValueSet[*Change](
		ChangeKeyer{},
	)
	c.addedNamed = collections_value.MakeMutableValueSet[*Change](
		ChangeKeyer{},
	)
	c.allB = collections_value.MakeMutableValueSet[values.String](nil)

	for h, es1 := range c.compareMapB.Named {
		var (
			existsInA bool
			es2       *metadatei.Metadatei
		)

		if es2, existsInA = c.compareMapA.Named[h]; existsInA {
			if es2.Bezeichnung.String() != es1.Bezeichnung.String() {
				c.modified.Add(&ChangeBezeichnung{Kennung: h, Bezeichnung: es1.Bezeichnung})
			}
		}

		change := Change{
			Key:     h,
			added:   kennung.MakeEtikettMutableSet(),
			removed: kennung.MakeEtikettMutableSet(),
		}

		c.allB.Add(values.MakeString(h))

		if err = es1.GetEtiketten().Each(
			func(e1 kennung.Etikett) (err error) {
				if c.compareMapA.Named.ContainsEtikett(h, e1) {
					// zettel had etikett previously
				} else {
					change.added.Add(e1)
				}

				if es2, ok := c.compareMapA.Named[h]; ok {
					es2.GetEtikettenMutable().Del(e1)
					c.compareMapA.Named[h] = es2
				}

				return
			},
		); err != nil {
			return
		}

		if existsInA {
			c.existing.Add(&change)
		} else {
			c.addedNamed.Add(&change)
		}
	}

	for h, es := range c.compareMapA.Named {
		var change *Change
		ok := false

		if change, ok = c.existing.Get(h); !ok {
			change = &Change{
				Key:     h,
				added:   kennung.MakeEtikettMutableSet(),
				removed: kennung.MakeEtikettMutableSet(),
			}
		}

		if err = es.GetEtiketten().Each(
			func(e1 kennung.Etikett) (err error) {
				if e1.String() == "" {
					err = errors.Errorf("empty etikett for %s", h)
					return
				}

				change.removed.Add(e1)

				return
			},
		); err != nil {
			return
		}

		c.existing.Add(change)
	}

	for h, es := range c.compareMapB.Unnamed {
		change := Change{
			Key:     h,
			added:   kennung.MakeEtikettMutableSet(),
			removed: kennung.MakeEtikettMutableSet(),
		}

		es.GetEtiketten().Each(change.added.Add)
		c.addedUnnamed.Add(&change)
	}

	return
}
