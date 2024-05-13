package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type constructor struct {
	Text
	all, implicit sku.SetPrefixVerzeichnisse
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
	c.all = sku.MakeSetPrefixVerzeichnisse(0)
	c.implicit = sku.MakeSetPrefixVerzeichnisse(0)
	mes := kennung.MakeMutableEtikettSet()

	if err = c.Transacted.Each(
		func(sk *sku.Transacted) (err error) {
			if err = c.all.Add(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			// eis := kennung.ExpandMany(
			// 	sk.Metadatei.Verzeichnisse.GetImplicitEtiketten(),
			// 	expansion.ExpanderRight,
			// )

			for _, ep := range sk.Metadatei.Verzeichnisse.Etiketten {
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

				if err = mes.Add(e); err != nil {
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

	if mes.Len() == 0 {
		c.EtikettSet = c.rootEtiketten
		return
	}

	if err = c.rootEtiketten.Each(mes.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.ExtraEtiketten.Each(mes.Add); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.ExtraEtiketten = mes
	c.EtikettSet = kennung.MakeEtikettSet()

	return
}

func (c *constructor) populate() (err error) {
	for _, e := range iter.Elements(c.ExtraEtiketten) {
		ee := newAssignment(c.GetDepth() + 1)
		ee.Etiketten = kennung.MakeEtikettSet(e)
		c.addChild(ee)

		segments := c.all.Subset(e)

		var used sku.TransactedMutableSet

		// ui.Debug().Print(e, segments)

		if used, err = c.makeChildren(
			ee,
			segments.Grouped,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		c.all = c.all.Subtract(used)
	}

	if _, err = c.makeChildren(
		c.Assignment,
		c.all,
	); err != nil {
		err = errors.Wrapf(err, "Assignment: %#v", c.Assignment)
		return
	}

	return
}

func (c *constructor) makeChildren(
	parent *Assignment,
	prefixSet sku.SetPrefixVerzeichnisse,
) (used sku.TransactedMutableSet, err error) {
	if c.GroupingEtiketten.Len() == 0 {
		if used, err = c.makeChildrenWithoutGroups(
			parent,
			prefixSet,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if used, err = c.makeChildrenWithGroups(
			parent,
			prefixSet,
			c.GroupingEtiketten,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c *constructor) makeChildrenWithoutGroups(
	parent *Assignment,
	prefixSet sku.SetPrefixVerzeichnisse,
) (used sku.TransactedMutableSet, err error) {
	used = sku.MakeTransactedMutableSet()
	prefixSet.EachZettel(used.Add)

	if err = c.makeAndAddUngrouped(parent, prefixSet.EachZettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *constructor) makeChildrenWithGroups(
	parent *Assignment,
	prefixSet sku.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.EtikettSlice,
) (used sku.TransactedMutableSet, err error) {
	if groupingEtiketten.Len() == 0 {
		if used, err = c.makeChildrenWithoutGroups(
			parent,
			prefixSet,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	used = sku.MakeTransactedMutableSet()
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
	segments sku.SetPrefixVerzeichnisseSegments,
	groupingEtiketten kennung.EtikettSlice,
	used sku.TransactedMutableSet,
) (err error) {
	if c.UsePrefixJoints {
		if err = c.addGroupedChildrenWithPrefixJoints(
			parent,
			segments.Grouped,
			groupingEtiketten,
			used,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = c.addGroupedChildrenWithoutPrefixJoints(
			parent,
			segments.Grouped,
			groupingEtiketten,
			used,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c *constructor) addGroupedChildrenWithoutPrefixJoints(
	parent *Assignment,
	prefixSet sku.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.EtikettSlice,
	used sku.TransactedMutableSet,
) (err error) {
	if err = prefixSet.Each(
		func(e kennung.Etikett, zs sku.TransactedMutableSet) (err error) {
			child := newAssignment(parent.GetDepth() + 1)
			child.Etiketten = kennung.MakeEtikettSet(e)

			nextGroupingEtiketten := kennung.MakeEtikettSlice()

			if groupingEtiketten.Len() > 1 {
				nextGroupingEtiketten = groupingEtiketten[1:]
			}

			var usedChild sku.TransactedMutableSet

			psv := sku.MakeSetPrefixVerzeichnisseFrom(zs)

			if usedChild, err = c.makeChildrenWithGroups(
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

func (c *constructor) addGroupedChildrenWithPrefixJoints(
	parent *Assignment,
	prefixSet sku.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.EtikettSlice,
	used sku.TransactedMutableSet,
) (err error) {
	if err = prefixSet.Each(
		func(e kennung.Etikett, zs sku.TransactedMutableSet) (err error) {
			if parent.Etiketten.Len() > 1 {
				return
			}

			prefixJoint := kennung.MakeEtikettSet(groupingEtiketten[0])

			var intermediate, lastChild *Assignment

			if len(parent.Children) > 0 {
				lastChild = parent.Children[len(parent.Children)-1]
			}

			if lastChild != nil &&
				(iter.SetEqualsPtr(lastChild.Etiketten, prefixJoint) ||
					lastChild.Etiketten.Len() == 0) {
				intermediate = lastChild
			} else {
				intermediate = newAssignment(parent.GetDepth() + 1)
				intermediate.Etiketten = prefixJoint
				parent.addChild(intermediate)
			}

			child := newAssignment(intermediate.GetDepth() + 1)

			var ls kennung.Etikett
			b := groupingEtiketten[0]

			if e.Equals(b) {
				return
			}

			if ls, err = kennung.LeftSubtract(e, b); err != nil {
				err = errors.Wrap(err)
				return
			}

			child.Etiketten = kennung.MakeEtikettSet(ls)

			nextGroupingEtiketten := kennung.MakeEtikettSlice()

			if groupingEtiketten.Len() > 1 {
				nextGroupingEtiketten = groupingEtiketten[1:]
			}

			var usedChild sku.TransactedMutableSet

			psv := sku.MakeSetPrefixVerzeichnisseFrom(zs)

			if usedChild, err = c.makeChildrenWithGroups(
				child,
				psv,
				nextGroupingEtiketten,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			usedChild.Each(used.Add)

			intermediate.addChild(child)

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
