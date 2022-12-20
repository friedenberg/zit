package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	zettel_pkg "github.com/friedenberg/zit/src/kilo/zettel"
)

type AssignmentTreeConstructor struct {
	Options
}

func (atc *AssignmentTreeConstructor) Assignments() (roots []*assignment, err error) {
	roots = make([]*assignment, 0, 1+atc.ExtraEtiketten.Len())

	root := newAssignment(0)
	root.etiketten = atc.RootEtiketten
	roots = append(roots, root)

	prefixSet := atc.Transacted.ToSetPrefixVerzeichnisse()

	for _, e := range atc.ExtraEtiketten.Elements() {
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
	prefixSet zettel_pkg.SetPrefixVerzeichnisse,
	groupingEtiketten kennung.Slice,
) (err error) {
	if groupingEtiketten.Len() == 0 {
		err = prefixSet.EachZettel(
			func(e kennung.Etikett, tz zettel_pkg.Verzeichnisse) (err error) {
				var z zettel

				if z, err = makeZettel(&tz, atc.Abbr); err != nil {
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
		func(tz *zettel_pkg.Verzeichnisse) (err error) {
			var z zettel

			if z, err = makeZettel(tz, atc.Abbr); err != nil {
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
		func(e kennung.Etikett, zs zettel_pkg.MutableSet) (err error) {
			if atc.UsePrefixJoints {
				if parent.etiketten.Len() > 1 {
				} else {
					prefixJoint := kennung.MakeEtikettSet(groupingEtiketten[0])

					var intermediate, lastChild *assignment

					if len(parent.children) > 0 {
						lastChild = parent.children[len(parent.children)-1]
					}

					if lastChild != nil && lastChild.etiketten.Equals(prefixJoint) {
						intermediate = lastChild
					} else {
						intermediate = newAssignment(parent.Depth() + 1)
						intermediate.etiketten = prefixJoint
						parent.addChild(intermediate)
					}

					child := newAssignment(intermediate.Depth() + 1)

					var ls kennung.Etikett

					if ls, err = e.LeftSubtract(groupingEtiketten[0]); err != nil {
						err = errors.Wrap(err)
						return
					}

					child.etiketten = kennung.MakeEtikettSet(ls)

					nextGroupingEtiketten := kennung.MakeSlice()

					if groupingEtiketten.Len() > 1 {
						nextGroupingEtiketten = groupingEtiketten[1:]
					}

					err = atc.makeChildren(child, zs.ToSetPrefixVerzeichnisse(), nextGroupingEtiketten)

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

				err = atc.makeChildren(child, zs.ToSetPrefixVerzeichnisse(), nextGroupingEtiketten)

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
		return parent.children[i].etiketten.String() < parent.children[j].etiketten.String()
	})

	return
}
