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
}

func (atc AssignmentTreeConstructor) RootAssignment() (root *assignment) {
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
	remainingEtiketten etikett.Slice,
) (assigned *stored_zettel.SetNamed) {
	assigned = stored_zettel.NewSetNamed()

	// logz.Print("making children")
	if remainingEtiketten.Len() == 0 {
		for _, zs := range prefixSet {
			// assigned.Merge(zs)
			for _, z := range zs {
				// logz.Printf("%s adding named %s", parent.etiketten, z.Hinweis)
				parent.named.Add(makeZettel(z))
			}
		}

		return
	}

	segments := prefixSet.Subset(remainingEtiketten[0])
	// logz.Printf("head: %s ungrouped: %s", remainingEtiketten[0], segments.Ungrouped.HinweisStrings())
	// logz.Printf("head: %s grouped: %s", remainingEtiketten[0], segments.Grouped.ToSetNamed().HinweisStrings())

	for _, z := range *segments.Ungrouped {
		parent.named.Add(makeZettel(z))
	}

	for e, zs := range *segments.Grouped {
		// assigned.Merge(zs)
		// logz.Print("iterating through grouped: ", e)
		child := newAssignment(parent.depth + 1)
		child.etiketten = etikett.MakeSet(e)
		// child.named = makeZettelZetFromSetNamed(zs)

		nextEtiketten := etikett.NewSlice()

		if remainingEtiketten.Len() > 1 {
			nextEtiketten = remainingEtiketten[1:]
		}

		_ = atc.makeChildren(child, *zs.ToSetPrefixNamed(), nextEtiketten)
		// childAssigned.Merge(c)
		// assigned.Merge(*c)

		parent.addChild(child)

		sort.Slice(parent.children, func(i, j int) bool {
			return parent.children[i].etiketten.String() < parent.children[j].etiketten.String()
		})
	}

	return
}
