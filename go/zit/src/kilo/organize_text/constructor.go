package organize_text

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/etiketten_path"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
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

	if err = c.Transacted.Each(c.all.AddTransacted); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.preparePrefixSetsAndRootsAndExtras(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// c.EtikettSet = c.rootEtiketten

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
) (explicitCount, implicitCount int, err error) {
	res := catgut.MakeFromString(re.String())

	if err = skus.Each(
		func(sk *sku.Transacted) (err error) {
			for _, ewp := range sk.Metadatei.Verzeichnisse.Etiketten.All {
				if ewp.Etikett.String() == sk.Kennung.String() {
					continue
				}

				cmp := ewp.ComparePartial(res)

				if cmp != 0 {
					continue
				}

				if len(ewp.Parents) == 0 { // TODO use Type
					explicitCount++
					break
				}

				for _, p := range ewp.Parents {
					if p.Type == etiketten_path.TypeDirect {
						explicitCount++
					} else {
						implicitCount++
					}
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
	anchored := kennung.MakeMutableEtikettSet()
	extras := kennung.MakeMutableEtikettSet()

	if err = c.rootEtiketten.Each(
		func(re kennung.Etikett) (err error) {
			var explicitCount, implicitCount int

			if explicitCount, implicitCount, err = c.collectExplicitAndImplicitFor(
				c.Transacted,
				re,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			ui.Log().Print(re, "explicit", explicitCount, "implicit", implicitCount)

			if explicitCount == c.Transacted.Len() {
				if err = anchored.Add(re); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else if explicitCount > 0 {
				if err = extras.Add(re); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.EtikettSet = anchored
	c.ExtraEtiketten = extras

	// c.ExtraEtiketten = implicit
	return
}

func (c *constructor) populate() (err error) {
	allUsed := makeObjSet()

	for _, e := range iter.Elements(c.ExtraEtiketten) {
		ee := c.makeChild(e)

		segments := c.all.Subset(e)

		if err = c.makeChildrenWithPossibleGroups(
			ee,
			segments.Grouped,
			c.GroupingEtiketten,
			allUsed,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

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

	if err = c.makeChildrenWithPossibleGroups(
		c.Assignment,
		c.all,
		c.GroupingEtiketten,
		allUsed,
	); err != nil {
		err = errors.Wrapf(err, "Assignment: %#v", c.Assignment)
		return
	}

	return
}

func (c *constructor) makeChildrenWithoutGroups(
	parent *Assignment,
	fi func(schnittstellen.FuncIter[*obj]) error,
	used objSet,
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
	used objSet,
) (err error) {
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

	segments := prefixSet.Subset(groupingEtiketten[0])

	if err = c.makeAndAddUngrouped(parent, segments.Ungrouped.Each); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.addGroupedChildren(
		parent,
		segments.Grouped,
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
	grouped PrefixSet,
	groupingEtiketten kennung.EtikettSlice,
	used objSet,
) (err error) {
	if err = grouped.Each(
		func(e kennung.Etikett, zs objSet) (err error) {
			if e.IsEmpty() || c.EtikettSet.Contains(e) {
				if err = c.makeAndAddUngrouped(parent, zs.Each); err != nil {
					err = errors.Wrap(err)
					return
				}

				zs.Each(used.Add)

				return
			}

			child := newAssignment(parent.GetDepth() + 1)
			child.Etiketten = kennung.MakeEtikettSet(e)
			groupingEtiketten.DropFirst()

			psv := MakePrefixSetFrom(zs)

			if err = c.makeChildrenWithPossibleGroups(
				child,
				psv,
				groupingEtiketten,
				used,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

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
	fi func(schnittstellen.FuncIter[*obj]) error,
) (err error) {
	if err = fi(
		func(tz *obj) (err error) {
			var z *obj

			if z, err = c.cloneObj(tz); err != nil {
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

func (c *constructor) cloneObj(
	named *obj,
) (z *obj, err error) {
	errors.TodoP4("add bez in a better way")

	z = &obj{Type: named.Type}

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
