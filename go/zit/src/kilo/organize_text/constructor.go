package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type constructor struct {
	Text
	all PrefixSet
}

func (c *constructor) Make() (ot *Text, err error) {
	ot = &c.Text
	c.all = MakePrefixSet(0)
	c.Assignment = newAssignment(0)
	c.IsRoot = true

	if err = c.Transacted.Each(c.all.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.preparePrefixSetsAndRootsAndExtras(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.populate(); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Metadatei.Typ = c.Options.Typ

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ot.SortChildren()

	return
}

func (c *constructor) collectExplicitAndImplicitFor(
	skus sku.TransactedSet,
	re kennung.Etikett,
	explicit kennung.EtikettMutableSet,
	implicit kennung.EtikettMutableSet,
	f schnittstellen.FuncIter[*sku.Transacted],
) (explicitCount, implicitCount int, err error) {
	res := catgut.MakeFromString(re.String())

	if err = skus.Each(
		func(sk *sku.Transacted) (err error) {
			if f != nil {
				if err = f(sk); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			for _, ewp := range sk.Metadatei.Verzeichnisse.Etiketten.All {
				if cmp := ewp.ComparePartial(res); cmp != 0 {
					continue
				}

				if len(ewp.Parents) > 0 {
					for _, p := range ewp.Parents {
						implicitCount++
						implicit.Add(kennung.MustEtikett(p.First().String()))
					}
				} else {
					explicitCount++
					explicit.Add(kennung.MustEtikett(ewp.Etikett.String()))
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *constructor) preparePrefixSetsAndRootsAndExtras() (err error) {
	implicit := kennung.MakeMutableEtikettSet()
	explicit := kennung.MakeMutableEtikettSet()

	mes := kennung.MakeMutableEtikettSet()

	if err = c.rootEtiketten.Each(
		func(re kennung.Etikett) (err error) {
			explicitCount := 0

			if explicitCount, _, err = c.collectExplicitAndImplicitFor(
				c.Transacted,
				re,
				explicit,
				implicit,
				nil,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if explicitCount != c.Transacted.Len() {
				return
			}

			if err = mes.Add(re); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.EtikettSet = mes
	// c.ExtraEtiketten = implicit

	return
}

func (c *constructor) populate() (err error) {
	allUsed := sku.MakeTransactedMutableSet()

	for _, e := range iter.Elements(c.ExtraEtiketten) {
		ee := c.makeChild(e)

		segments := c.all.Match(e)

		var used sku.TransactedMutableSet

		if used, err = c.makeChildrenWithPossibleGroups(
			ee,
			segments.Grouped,
			c.GroupingEtiketten,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		used.Each(allUsed.Add)

		if err = c.makeChildrenWithoutGroups(
			ee,
			segments.Ungrouped.Each,
			allUsed,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	c.all = c.all.Subtract(allUsed)

	if _, err = c.makeChildrenWithPossibleGroups(
		c.Assignment,
		c.all,
		c.GroupingEtiketten,
	); err != nil {
		err = errors.Wrapf(err, "Assignment: %#v", c.Assignment)
		return
	}

	return
}

func (c *constructor) makeChildrenWithoutGroups(
	parent *Assignment,
	fi func(schnittstellen.FuncIter[*sku.Transacted]) error,
	used sku.TransactedMutableSet,
) (err error) {
	if err = fi(used.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.makeAndAddUngrouped(parent, fi); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *constructor) makeChildrenWithPossibleGroups(
	parent *Assignment,
	prefixSet PrefixSet,
	groupingEtiketten kennung.EtikettSlice,
) (used sku.TransactedMutableSet, err error) {
	used = sku.MakeTransactedMutableSet()

	if groupingEtiketten.Len() == 0 {
		if err = c.makeChildrenWithoutGroups(
			parent,
			prefixSet.EachZettel,
			used,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	// TODO use implicit here
	segments := prefixSet.Subset(groupingEtiketten[0])

	if err = c.makeAndAddUngrouped(parent, segments.Ungrouped.Each); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.addGroupedChildren(
		parent,
		segments,
		groupingEtiketten,
		used,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	parent.SortChildren()

	return
}

func (c *constructor) addGroupedChildren(
	parent *Assignment,
	segments Segments,
	groupingEtiketten kennung.EtikettSlice,
	used sku.TransactedMutableSet,
) (err error) {
	if err = segments.Grouped.Each(
		func(e kennung.Etikett, zs sku.TransactedMutableSet) (err error) {
			child := newAssignment(parent.GetDepth() + 1)
			child.Etiketten = kennung.MakeEtikettSet(e)

			nextGroupingEtiketten := kennung.MakeEtikettSlice()

			if groupingEtiketten.Len() > 1 {
				nextGroupingEtiketten = groupingEtiketten[1:]
			}

			var usedChild sku.TransactedMutableSet

			psv := MakePrefixSetFrom(zs)

			if usedChild, err = c.makeChildrenWithPossibleGroups(
				child,
				psv,
				nextGroupingEtiketten,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			usedChild.Each(used.Add)

			parent.addChild(child)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *constructor) makeAndAddUngrouped(
	parent *Assignment,
	fi func(schnittstellen.FuncIter[*sku.Transacted]) error,
) (err error) {
	if err = fi(
		func(tz *sku.Transacted) (err error) {
			var z *obj

			if z, err = c.makeObj(tz); err != nil {
				err = errors.Wrap(err)
				return
			}

			parent.AddObjekte(z)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}
	return
}

func (c *constructor) makeObj(
	named *sku.Transacted,
) (z *obj, err error) {
	errors.TodoP4("add bez in a better way")

	z = &obj{}

	if err = z.SetFromSkuLike(named); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.removeEtikettenIfNecessary(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *constructor) removeEtikettenIfNecessary(
	o *obj,
) (err error) {
	if c.PrintOptions.PrintEtikettenAlways {
		return
	}

	if o.Metadatei.Bezeichnung.IsEmpty() {
		return
	}

	o.Metadatei.ResetEtiketten()

	return
}
