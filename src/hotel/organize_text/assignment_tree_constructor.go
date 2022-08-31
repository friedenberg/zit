package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/hotel/collections"
)

type AssignmentTreeConstructor struct {
	RootEtiketten     etikett.Set
	GroupingEtiketten etikett.Slice
	ExtraEtiketten    etikett.Set
	Transacted        collections.SetTransacted
	UsePrefixJoints   bool
}

func (atc *AssignmentTreeConstructor) Assignments() (roots []*assignment) {
	roots = make([]*assignment, 0, 1+len(atc.ExtraEtiketten))

	root := newAssignment()
	root.etiketten = atc.RootEtiketten
	roots = append(roots, root)

	prefixSet := atc.Transacted.ToSetPrefixTransacted()
	atc.makeChildren(root, prefixSet, atc.GroupingEtiketten)

	for _, e := range atc.ExtraEtiketten {
		child := newAssignment()
		child.etiketten = etikett.MakeSet(e)
		roots = append(roots, child)
	}

	return
}

func (atc AssignmentTreeConstructor) makeChildren(
	parent *assignment,
	prefixSet collections.SetPrefixTransacted,
	groupingEtiketten etikett.Slice,
) {
	if groupingEtiketten.Len() == 0 {
		prefixSet.EachZettel(
			func(e etikett.Etikett, tz zettel_stored.Transacted) (err error) {
				parent.named.Add(makeZettel(tz.Named))

				return
			},
		)

		return
	}

	segments := prefixSet.Subset(groupingEtiketten[0])

	segments.Ungrouped.Each(
		func(tz zettel_stored.Transacted) (err error) {
			parent.named.Add(makeZettel(tz.Named))
			return
		},
	)

	segments.Grouped.Each(
		func(e etikett.Etikett, zs collections.SetTransacted) (err error) {
			if atc.UsePrefixJoints {
				if parent.etiketten.Len() > 1 {
				} else {
					prefixJoint := etikett.MakeSet(groupingEtiketten[0])

					var intermediate, lastChild *assignment

					if len(parent.children) > 0 {
						lastChild = parent.children[len(parent.children)-1]
					}

					if lastChild != nil && lastChild.etiketten.Equals(prefixJoint) {
						intermediate = lastChild
					} else {
						intermediate = newAssignment()
						intermediate.etiketten = prefixJoint
						parent.addChild(intermediate)
					}

					child := newAssignment()
					child.etiketten = etikett.MakeSet(e.LeftSubtract(groupingEtiketten[0]))

					nextGroupingEtiketten := etikett.NewSlice()

					if groupingEtiketten.Len() > 1 {
						nextGroupingEtiketten = groupingEtiketten[1:]
					}

					atc.makeChildren(child, zs.ToSetPrefixTransacted(), nextGroupingEtiketten)

					intermediate.addChild(child)
				}
			} else {
				child := newAssignment()
				child.etiketten = etikett.MakeSet(e)

				nextGroupingEtiketten := etikett.NewSlice()

				if groupingEtiketten.Len() > 1 {
					nextGroupingEtiketten = groupingEtiketten[1:]
				}

				atc.makeChildren(child, zs.ToSetPrefixTransacted(), nextGroupingEtiketten)

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
