package organize_text

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
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

	c.Metadata.Typ = c.Options.Type

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ot.SortChildren()

	return
}

func (c *constructor) collectExplicitAndImplicitFor(
	skus sku.TransactedSet,
	re ids.Tag,
) (explicitCount, implicitCount int, err error) {
	res := catgut.MakeFromString(re.String())

	if err = skus.Each(
		func(sk *sku.Transacted) (err error) {
			for _, ewp := range sk.Metadata.Cache.TagPaths.All {
				if ewp.Etikett.String() == sk.ObjectId.String() {
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
					if p.Type == tag_paths.TypeDirect {
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
	anchored := ids.MakeMutableTagSet()
	extras := ids.MakeMutableTagSet()

	if err = c.rootTags.Each(
		func(re ids.Tag) (err error) {
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

	c.TagSet = anchored
	c.ExtraTags = extras

	// c.ExtraEtiketten = implicit
	return
}

func (c *constructor) populate() (err error) {
	allUsed := makeObjSet()

	for _, e := range iter.Elements(c.ExtraTags) {
		ee := c.makeChild(e)

		segments := c.all.Subset(e)

		if err = c.makeChildrenWithPossibleGroups(
			ee,
			segments.Grouped,
			c.GroupingTags,
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
		c.GroupingTags,
		allUsed,
	); err != nil {
		err = errors.Wrapf(err, "Assignment: %#v", c.Assignment)
		return
	}

	return
}

func (c *constructor) makeChildrenWithoutGroups(
	parent *Assignment,
	fi func(interfaces.FuncIter[*obj]) error,
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
	groupingEtiketten ids.TagSlice,
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
	groupingEtiketten ids.TagSlice,
	used objSet,
) (err error) {
	if err = grouped.Each(
		func(e ids.Tag, zs objSet) (err error) {
			if e.IsEmpty() || c.TagSet.Contains(e) {
				if err = c.makeAndAddUngrouped(parent, zs.Each); err != nil {
					err = errors.Wrap(err)
					return
				}

				zs.Each(used.Add)

				return
			}

			child := newAssignment(parent.GetDepth() + 1)
			child.Tags = ids.MakeTagSet(e)
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
	fi func(interfaces.FuncIter[*obj]) error,
) (err error) {
	if err = fi(
		func(tz *obj) (err error) {
			var z *obj

			if z, err = c.cloneObj(tz); err != nil {
				err = errors.Wrap(err)
				return
			}

			parent.AddObject(z)

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

	if o.Metadata.Description.IsEmpty() {
		return
	}

	o.Metadata.ResetTags()

	return
}
