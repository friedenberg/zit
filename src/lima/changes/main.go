package changes

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/kilo/organize_text"
)

type Changes interface {
	GetModified() schnittstellen.SetLike[ChangeBezeichnung]
	GetExisting() schnittstellen.SetLike[Change]
	GetAddedUnnamed() schnittstellen.SetLike[Change]
	GetAddedNamed() schnittstellen.SetLike[Change]
	GetAllBKeys() schnittstellen.SetLike[values.String]
}

type changes struct {
	modified     schnittstellen.MutableSetLike[ChangeBezeichnung]
	existing     schnittstellen.MutableSetLike[Change]
	addedUnnamed schnittstellen.MutableSetLike[Change]
	addedNamed   schnittstellen.MutableSetLike[Change]
	allB         schnittstellen.MutableSetLike[values.String]
}

func (c changes) GetModified() schnittstellen.SetLike[ChangeBezeichnung] {
	return c.modified
}

func (c changes) GetExisting() schnittstellen.SetLike[Change] {
	return c.existing
}

func (c changes) GetAddedUnnamed() schnittstellen.SetLike[Change] {
	return c.addedUnnamed
}

func (c changes) GetAddedNamed() schnittstellen.SetLike[Change] {
	return c.addedNamed
}

func (c changes) GetAllBKeys() schnittstellen.SetLike[values.String] {
	return c.allB
}

func makeCompareMapFromOrganizeTextAndExpander(
	in *organize_text.Text,
	hinweis_expander func(string) (kennung.Hinweis, error),
) (out organize_text.CompareMap, err error) {
	var preExpansion organize_text.CompareMap

	if preExpansion, err = in.ToCompareMap(); err != nil {
		err = errors.Wrap(err)
		return
	}

	out = organize_text.CompareMap{
		Named:   make(organize_text.SetKeyToMetadatei),
		Unnamed: preExpansion.Unnamed,
	}

	for h, v := range preExpansion.Named {
		var h1 schnittstellen.Stringer

		if h1, err = hinweis_expander(h); err == nil {
			h = h1.String()
		}

		err = nil

		out.Named[h] = v
	}

	return
}

func ChangesFrom(
	a1, b1 *organize_text.Text,
	hinweis_expander func(string) (kennung.Hinweis, error),
) (c1 Changes, err error) {
	var c changes
	c1 = &c
	var a, b organize_text.CompareMap

	if a, err = makeCompareMapFromOrganizeTextAndExpander(
		a1,
		hinweis_expander,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if b, err = makeCompareMapFromOrganizeTextAndExpander(
		b1,
		hinweis_expander,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.modified = collections_ptr.MakeMutableSet[ChangeBezeichnung, *ChangeBezeichnung](
		ChangeBezeichnungKeyer{},
	)
	c.existing = collections_ptr.MakeMutableValueSet[Change, *Change](
		ChangeKeyer{},
	)
	c.addedUnnamed = collections_ptr.MakeMutableValueSet[Change, *Change](
		ChangeKeyer{},
	)
	c.addedNamed = collections_ptr.MakeMutableValueSet[Change, *Change](
		ChangeKeyer{},
	)
	c.allB = collections_value.MakeMutableValueSet[values.String](nil)

	for h, es1 := range b.Named {
		var (
			existsInA bool
			es2       metadatei.Metadatei
		)

		if es2, existsInA = a.Named[h]; existsInA {
			if es2.Bezeichnung.String() != es1.Bezeichnung.String() {
				c.modified.Add(ChangeBezeichnung{Kennung: h, Bezeichnung: es1.Bezeichnung})
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
				if a.Named.ContainsEtikett(h, e1) {
					// zettel had etikett previously
				} else {
					change.added.Add(e1)
				}

				if es2, ok := a.Named[h]; ok {
					es2.GetEtikettenMutable().Del(e1)
					a.Named[h] = es2
				}

				return
			},
		); err != nil {
			return
		}

		if existsInA {
			c.existing.Add(change)
		} else {
			c.addedNamed.Add(change)
		}
	}

	for h, es := range a.Named {
		var change Change
		ok := false

		if change, ok = c.existing.Get(h); !ok {
			change = Change{
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

	for h, es := range b.Unnamed {
		change := Change{
			Key:     h,
			added:   kennung.MakeEtikettMutableSet(),
			removed: kennung.MakeEtikettMutableSet(),
		}

		es.GetEtiketten().Each(change.added.Add)
		c.addedUnnamed.Add(change)
	}

	return
}
