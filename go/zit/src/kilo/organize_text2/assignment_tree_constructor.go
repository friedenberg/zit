package organize_text2

import (
	"sort"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/objekte_collections"
)

type AssignmentTreeConstructor struct {
	Options
}

func (atc *AssignmentTreeConstructor) Assignments() (roots []*Assignment, err error) {
	roots = make([]*Assignment, 0, 1+atc.ExtraEtiketten.Len())

	root := newAssignment(0)
	root.Etiketten = atc.rootEtiketten
	roots = append(roots, root)

	prefixSet := objekte_collections.MakeSetPrefixVerzeichnisse(0)
	atc.Transacted.Each(prefixSet.Add)

	for _, e := range iter.Elements[kennung.Etikett](atc.ExtraEtiketten) {
		errors.Err().Printf("making extras: %s", e)
		errors.Err().Printf("prefix set before: %v", prefixSet)
		if err = atc.makeChildren(root, prefixSet, kennung.MakeEtikettSlice(e)); err != nil {
			err = errors.Wrap(err)
			return
		}
		errors.Err().Printf("prefix set after: %v", prefixSet)
	}

	if err = atc.makeChildren(root, prefixSet, atc.GroupingEtiketten); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (atc AssignmentTreeConstructor) makeChildren(
	parent *Assignment,
	prefixSet objekte_collections.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.EtikettSlice,
) (err error) {
	if groupingEtiketten.Len() == 0 {
		err = prefixSet.EachZettel(
			func(e kennung.Etikett, tz *sku.Transacted) (err error) {
				var z *obj

				if z, err = makeObj(
					atc.PrintOptions,
					tz,
					atc.Expanders,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				parent.Named.Add(z)

				return
			},
		)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	segments := prefixSet.Subset(groupingEtiketten[0])

	err = segments.Ungrouped.Each(
		func(tz *sku.Transacted) (err error) {
			var z *obj

			if z, err = makeObj(
				atc.PrintOptions,
				tz,
				atc.Expanders,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			parent.Named.Add(z)

			return
		},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	segments.Grouped.Each(
		func(e kennung.Etikett, zs objekte_collections.MutableSetMetadateiWithKennung) (err error) {
			if atc.UsePrefixJoints {
				if parent.Etiketten != nil && parent.Etiketten.Len() > 1 {
				} else {
					prefixJoint := kennung.MakeEtikettSet(groupingEtiketten[0])

					var intermediate, lastChild *Assignment

					if len(parent.Children) > 0 {
						lastChild = parent.Children[len(parent.Children)-1]
					}

					if lastChild != nil && iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](lastChild.Etiketten, prefixJoint) {
						intermediate = lastChild
					} else {
						intermediate = newAssignment(parent.GetDepth() + 1)
						intermediate.Etiketten = prefixJoint
						parent.addChild(intermediate)
					}

					child := newAssignment(intermediate.GetDepth() + 1)

					var ls kennung.Etikett

					if ls, err = kennung.LeftSubtract(e, groupingEtiketten[0]); err != nil {
						err = errors.Wrap(err)
						return
					}

					child.Etiketten = kennung.MakeEtikettSet(ls)

					nextGroupingEtiketten := kennung.MakeEtikettSlice()

					if groupingEtiketten.Len() > 1 {
						nextGroupingEtiketten = groupingEtiketten[1:]
					}

					psv := objekte_collections.MakeSetPrefixVerzeichnisse(0)
					zs.Each(psv.Add)
					err = atc.makeChildren(child, psv, nextGroupingEtiketten)

					if err != nil {
						err = errors.Wrap(err)
						return
					}

					intermediate.addChild(child)
				}
			} else {
				child := newAssignment(parent.GetDepth() + 1)
				child.Etiketten = kennung.MakeEtikettSet(e)

				nextGroupingEtiketten := kennung.MakeEtikettSlice()

				if groupingEtiketten.Len() > 1 {
					nextGroupingEtiketten = groupingEtiketten[1:]
				}

				psv := objekte_collections.MakeSetPrefixVerzeichnisse(0)
				zs.Each(psv.Add)
				err = atc.makeChildren(child, psv, nextGroupingEtiketten)

				if err != nil {
					err = errors.Wrap(err)
					return
				}

				parent.addChild(child)
			}
			return
		},
	)

	sort.Slice(parent.Children, func(i, j int) bool {
		vi := iter.StringCommaSeparated[kennung.Etikett](
			parent.Children[i].Etiketten,
		)
		vj := iter.StringCommaSeparated[kennung.Etikett](
			parent.Children[j].Etiketten,
		)
		return vi < vj
	})

	return
}
