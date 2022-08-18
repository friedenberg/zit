package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
)

type AssignmentTreeConstructor struct {
	RootEtiketten     etikett.Set
	GroupingEtiketten etikett.Slice
	ExtraEtiketten    etikett.Set
	Named             stored_zettel.SetNamed
	UsePrefixJoints   bool
}

func (atc *AssignmentTreeConstructor) RootAssignment() (root *assignment) {
	// atc.UsePrefixJoints = true
	root = newAssignment(1)
	root.etiketten = atc.RootEtiketten

	prefixSet := atc.Named.ToSetPrefixNamed()
	atc.makeChildren(root, *prefixSet, atc.GroupingEtiketten)

	for _, e := range atc.ExtraEtiketten {
		child := newAssignment(root.depth + 1)
		child.etiketten = etikett.MakeSet(e)
		root.addChild(child)
	}

	return
}

func (atc AssignmentTreeConstructor) makeChildren(
	parent *assignment,
	prefixSet stored_zettel.SetPrefixNamed,
	groupingEtiketten etikett.Slice,
) {
	if groupingEtiketten.Len() == 0 {
		for _, zs := range prefixSet {
			for _, z := range zs {
				parent.named.Add(makeZettel(z))
			}
		}

		return
	}

	segments := prefixSet.Subset(groupingEtiketten[0])

	for _, z := range *segments.Ungrouped {
		parent.named.Add(makeZettel(z))
	}

	for e, zs := range *segments.Grouped {
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
					intermediate = newAssignment(parent.depth + 1)
					intermediate.etiketten = prefixJoint
					parent.addChild(intermediate)
				}

				child := newAssignment(intermediate.depth + 1)
				child.etiketten = etikett.MakeSet(e.LeftSubtract(groupingEtiketten[0]))

				nextGroupingEtiketten := etikett.NewSlice()

				if groupingEtiketten.Len() > 1 {
					nextGroupingEtiketten = groupingEtiketten[1:]
				}

				atc.makeChildren(child, *zs.ToSetPrefixNamed(), nextGroupingEtiketten)

				intermediate.addChild(child)
			}
		} else {
			child := newAssignment(parent.depth + 1)
			child.etiketten = etikett.MakeSet(e)

			nextGroupingEtiketten := etikett.NewSlice()

			if groupingEtiketten.Len() > 1 {
				nextGroupingEtiketten = groupingEtiketten[1:]
			}

			atc.makeChildren(child, *zs.ToSetPrefixNamed(), nextGroupingEtiketten)

			parent.addChild(child)
		}
	}

	sort.Slice(parent.children, func(i, j int) bool {
		return parent.children[i].etiketten.String() < parent.children[j].etiketten.String()
	})

	return
}
