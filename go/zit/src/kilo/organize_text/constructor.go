package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type constructor struct {
	Text
	all, implicit PrefixSet
}

func (c *constructor) Make() (ot *Text, err error) {
	ot = &c.Text
	c.Assignment = newAssignment(0)
	c.IsRoot = true

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

func (c *constructor) preparePrefixSetsAndRootsAndExtras() (err error) {
	c.all = MakePrefixSet(0)
	c.implicit = MakePrefixSet(0)
	implicit := kennung.MakeMutableEtikettSet()
	explicit := kennung.MakeMutableEtikettSet()

	if err = c.Transacted.Each(
		func(sk *sku.Transacted) (err error) {
			if err = c.all.Add(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			isExplicit := false

			if err = c.rootEtiketten.Each(
				func(re kennung.Etikett) (err error) {
					if err = sk.Metadatei.GetEtiketten().Each(
						func(skE kennung.Etikett) (err error) {
							es := kennung.ExpandOne(
								&skE,
								expansion.ExpanderRight,
							)

							if es.Contains(re) {
								explicit.Add(skE)
								isExplicit = true
							}

							return
						},
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					return
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if isExplicit {
				return
			}

			for _, ep := range sk.Metadatei.Verzeichnisse.Etiketten.Paths {
				if ep.Len() == 1 {
					continue
				}

				if !c.rootEtiketten.ContainsKey(ep.Last().String()) {
					continue
				}

				if err = c.implicit.Add(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				var e kennung.Etikett

				if err = e.Set(ep.First().String()); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = implicit.Add(e); err != nil {
					err = errors.Wrap(err)
					return
				}

				// ui.Debug().Print(sk, (*etiketten_path.StringBackward)(ep))
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if implicit.Len() == 0 {
		c.EtikettSet = c.rootEtiketten
		return
	}

	if err = explicit.Each(implicit.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.ExtraEtiketten.Each(implicit.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.ExtraEtiketten = implicit
	c.EtikettSet = kennung.MakeEtikettSet()

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

	o.Metadatei.GetEtikettenMutable().Reset()

	return
}
