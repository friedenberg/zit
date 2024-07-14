package organize_text

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type Refiner struct {
	Enabled         bool
	UsePrefixJoints bool
}

func (atc *Refiner) shouldMergeAllChildrenIntoParent(a *Assignment) (ok bool) {
	switch {
	case a.Parent.IsRoot:
		fallthrough

	default:
		ok = false
	}

	return
}

func (atc *Refiner) shouldMergeIntoParent(a *Assignment) bool {
	ui.Log().Printf("checking node should merge: %s", a)

	if a.Parent == nil {
		ui.Log().Print("parent is nil")
		return false
	}

	if a.Parent.IsRoot {
		ui.Log().Print("parent is root")
		return false
	}

	if a.Tags.Len() == 1 && ids.IsEmpty(a.Tags.Any()) {
		ui.Log().Print("1 Etikett, and it's empty, merging")
		return true
	}

	if a.Tags.Len() == 0 {
		ui.Log().Print("etiketten length is 0, merging")
		return true
	}

	if a.Parent.Tags.Len() != 1 {
		ui.Log().Print("parent etiketten length is not 1")
		return false
	}

	if a.Tags.Len() != 1 {
		ui.Log().Print("etiketten length is not 1")
		return false
	}

	equal := iter.SetEqualsPtr(a.Tags, a.Parent.Tags)

	if !equal {
		ui.Log().Print("parent etiketten not equal")
		return false
	}

	if ids.IsDependentLeaf(a.Parent.Tags.Any()) {
		ui.Log().Print("is prefix joint")
		return false
	}

	if ids.IsDependentLeaf(a.Tags.Any()) {
		ui.Log().Print("is prefix joint")
		return false
	}

	return true
}

func (atc *Refiner) renameForPrefixJoint(a *Assignment) (err error) {
	if !atc.UsePrefixJoints {
		return
	}

	if a == nil {
		ui.Log().Printf("assignment is nil")
		return
	}

	if a.Parent == nil {
		ui.Log().Printf("parent is nil: %#v", a)
		return
	}

	if a.Parent.Tags.Len() == 0 {
		return
	}

	if a.Parent.Tags.Len() != 1 {
		return
	}

	if ids.IsDependentLeaf(a.Parent.Tags.Any()) {
		return
	}

	if ids.IsDependentLeaf(a.Tags.Any()) {
		return
	}

	if !ids.HasParentPrefix(a.Tags.Any(), a.Parent.Tags.Any()) {
		ui.Log().Print("parent is not prefix joint")
		return
	}

	aEtt := a.Tags.Any()
	pEtt := a.Parent.Tags.Any()

	if aEtt.Equals(pEtt) {
		ui.Log().Print("parent is is equal to child")
		return
	}

	var ls ids.Tag

	if ls, err = ids.LeftSubtract(aEtt, pEtt); err != nil {
		err = errors.Wrap(err)
		return
	}

	a.Tags = ids.MakeTagSet(ls)

	return
}

// passed-in assignment may be nil?
func (atc *Refiner) Refine(a *Assignment) (err error) {
	if !atc.Enabled {
		return
	}

	if !a.IsRoot {
		if atc.shouldMergeIntoParent(a) {
			ui.Log().Print("merging into parent")
			p := a.Parent

			if err = p.consume(a); err != nil {
				err = errors.Wrap(err)
				return
			}

			return atc.Refine(p)
		}
	}

	if err = atc.applyPrefixJoints(a); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = atc.renameForPrefixJoint(a); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, child := range a.Children {
		if err = atc.Refine(child); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = atc.applyPrefixJoints(a); err != nil {
		err = errors.Wrap(err)
		return
	}

	a.SortChildren()

	return
}

func (atc Refiner) applyPrefixJoints(a *Assignment) (err error) {
	if !atc.UsePrefixJoints {
		return
	}

	if a.Tags == nil || a.Tags.Len() == 0 {
		return
	}

	childPrefixes := atc.childPrefixes(a)

	if len(childPrefixes) == 0 {
		return
	}

	groupingPrefix := childPrefixes[0]

	var na *Assignment

	if a.Tags.Len() == 1 &&
		a.Tags.Any().Equals(groupingPrefix.Tag) {
		na = a
	} else {
		na = newAssignment(a.GetDepth() + 1)
		na.Tags = ids.MakeTagSet(groupingPrefix.Tag)
		a.addChild(na)
	}

	for _, c := range groupingPrefix.assignments {
		if c.Parent != na {
			if err = c.removeFromParent(); err != nil {
				err = errors.Wrap(err)
				return
			}

			na.addChild(c)
		}

		c.Tags = ids.SubtractPrefix(
			c.Tags,
			groupingPrefix.Tag,
		)
	}

	return
}

type etikettBag struct {
	ids.Tag
	assignments []*Assignment
}

func (a Refiner) childPrefixes(node *Assignment) (out []etikettBag) {
	m := make(map[string][]*Assignment)
	out = make([]etikettBag, 0, len(node.Children))

	if node.Tags.Len() == 0 {
		return
	}

	for _, c := range node.Children {
		expanded := ids.Expanded(c.Tags, expansion.ExpanderRight)

		expanded.Each(
			func(e ids.Tag) (err error) {
				if e.String() == "" {
					return
				}

				var n []*Assignment
				ok := false

				if n, ok = m[e.String()]; !ok {
					n = make([]*Assignment, 0)
				}

				n = append(n, c)

				m[e.String()] = n

				return
			},
		)
	}

	for e, n := range m {
		if len(n) > 1 {
			var e1 ids.Tag

			errors.PanicIfError(e1.Set(e))

			out = append(out, etikettBag{Tag: e1, assignments: n})
		}
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			if len(out[i].assignments) == len(out[j].assignments) {
				return len(
					out[i].Tag.String(),
				) > len(
					out[j].Tag.String(),
				)
			} else {
				return len(out[i].assignments) > len(out[j].assignments)
			}
		},
	)

	return
}
