package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/objekte_collections"
)

type AssignmentTreeConstructor struct {
	Options
}

func (atc *AssignmentTreeConstructor) Assignments() (roots []*assignment, err error) {
	roots = make([]*assignment, 0, 1+atc.ExtraEtiketten.Len())

	root := newAssignment(0)
	root.etiketten = atc.rootEtiketten
	roots = append(roots, root)

	prefixSet := objekte_collections.MakeSetPrefixVerzeichnisse(0)
	atc.Transacted.Each(prefixSet.Add)

	for _, e := range iter.Elements[kennung.Etikett](atc.ExtraEtiketten) {
		errors.Err().Printf("making extras: %s", e)
		errors.Err().Printf("prefix set before: %v", prefixSet)
		if err = atc.makeChildren(root, prefixSet, kennung.MakeSlice(e)); err != nil {
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
	parent *assignment,
	prefixSet objekte_collections.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.Slice,
) (err error) {
	if groupingEtiketten.Len() == 0 {
		err = prefixSet.EachZettel(
			func(e kennung.Etikett, tz *sku.Transacted) (err error) {
				var z *obj

				if z, err = makeObj(
          atc.Options.PrintOptions,
          tz,
          atc.Expanders,
        ); err != nil {
					err = errors.Wrap(err)
					return
				}

				parent.named.Add(z)

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
				atc.Options.PrintOptions,
				tz,
				atc.Expanders,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			parent.named.Add(z)

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
				if parent.etiketten != nil && parent.etiketten.Len() > 1 {
				} else {
					prefixJoint := kennung.MakeEtikettSet(groupingEtiketten[0])

					var intermediate, lastChild *assignment

					if len(parent.children) > 0 {
						lastChild = parent.children[len(parent.children)-1]
					}

					if lastChild != nil && iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](lastChild.etiketten, prefixJoint) {
						intermediate = lastChild
					} else {
						intermediate = newAssignment(parent.Depth() + 1)
						intermediate.etiketten = prefixJoint
						parent.addChild(intermediate)
					}

					child := newAssignment(intermediate.Depth() + 1)

					var ls kennung.Etikett

					if ls, err = kennung.LeftSubtract(e, groupingEtiketten[0]); err != nil {
						err = errors.Wrap(err)
						return
					}

					child.etiketten = kennung.MakeEtikettSet(ls)

					nextGroupingEtiketten := kennung.MakeSlice()

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
				child := newAssignment(parent.Depth() + 1)
				child.etiketten = kennung.MakeEtikettSet(e)

				nextGroupingEtiketten := kennung.MakeSlice()

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

	sort.Slice(parent.children, func(i, j int) bool {
		vi := iter.StringCommaSeparated[kennung.Etikett](
			parent.children[i].etiketten,
		)
		vj := iter.StringCommaSeparated[kennung.Etikett](
			parent.children[j].etiketten,
		)
		return vi < vj
	})

	return
}
