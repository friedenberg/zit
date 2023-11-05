package changes

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/kilo/organize_text"
)

type Changes interface {
	GetExisting() schnittstellen.SetLike[Change]
	GetAdded() schnittstellen.SetLike[Change]
	GetAllBKeys() schnittstellen.SetLike[values.String]
}

type changes struct {
	existing schnittstellen.MutableSetLike[Change]
	added    schnittstellen.MutableSetLike[Change]
	allB     schnittstellen.MutableSetLike[values.String]
}

func (c changes) GetExisting() schnittstellen.SetLike[Change] {
	return c.existing
}

func (c changes) GetAdded() schnittstellen.SetLike[Change] {
	return c.added
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
		if h1, err := hinweis_expander(h); err == nil {
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

	c.existing = collections_ptr.MakeMutableValueSet[Change, *Change](
		ChangeKeyer{},
	)
	c.added = collections_ptr.MakeMutableValueSet[Change, *Change](
		ChangeKeyer{},
	)
	c.allB = collections_value.MakeMutableValueSet[values.String](nil)

	for h, es1 := range b.Named {
		change := Change{
			Key:     h,
			added:   kennung.MakeEtikettMutableSet(),
			removed: kennung.MakeEtikettMutableSet(),
		}

		c.allB.Add(values.MakeString(h))

		for _, e1 := range es1.GetEtiketten().Elements() {
			if a.Named.ContainsEtikett(h, e1) {
				// zettel had etikett previously
			} else {
				change.added.Add(e1)
			}

			if es2, ok := a.Named[h]; ok {
				es2.GetEtikettenMutable().Del(e1)
				a.Named[h] = es2
			}
		}

		c.existing.Add(change)
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

		for _, e1 := range es.GetEtiketten().Elements() {
			if e1.String() == "" {
				err = errors.Errorf("empty etikett for %s", h)
				return
			}

			change.removed.Add(e1)
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
		c.added.Add(change)
	}

	return
}
