package organize_text

import (
	"sort"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
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
	errors.Log().Printf("checking node should merge: %s", a)

	if a.Parent == nil {
		errors.Log().Print("parent is nil")
		return false
	}

	if a.Parent.IsRoot {
		errors.Log().Print("parent is root")
		return false
	}

	if a.Etiketten.Len() == 1 && kennung.IsEmpty(a.Etiketten.Any()) {
		errors.Log().Print("1 Etikett, and it's empty, merging")
		return true
	}

	if a.Etiketten.Len() == 0 {
		errors.Log().Print("etiketten length is 0, merging")
		return true
	}

	if a.Parent.Etiketten.Len() != 1 {
		errors.Log().Print("parent etiketten length is not 1")
		return false
	}

	if a.Etiketten.Len() != 1 {
		errors.Log().Print("etiketten length is not 1")
		return false
	}

	equal := iter.SetEqualsPtr(a.Etiketten, a.Parent.Etiketten)

	if !equal {
		errors.Log().Print("parent etiketten not equal")
		return false
	}

	if kennung.IsDependentLeaf(a.Parent.Etiketten.Any()) {
		errors.Log().Print("is prefix joint")
		return false
	}

	if kennung.IsDependentLeaf(a.Etiketten.Any()) {
		errors.Log().Print("is prefix joint")
		return false
	}

	return true
}

func (atc *Refiner) renameForPrefixJoint(a *Assignment) (err error) {
	if !atc.UsePrefixJoints {
		return
	}

	if a == nil {
		errors.Log().Printf("assignment is nil")
		return
	}

	if a.Parent == nil {
		errors.Log().Printf("parent is nil: %#v", a)
		return
	}

	if a.Parent.Etiketten.Len() == 0 {
		return
	}

	if a.Parent.Etiketten.Len() != 1 {
		return
	}

	if kennung.IsDependentLeaf(a.Parent.Etiketten.Any()) {
		return
	}

	if kennung.IsDependentLeaf(a.Etiketten.Any()) {
		return
	}

	if !kennung.HasParentPrefix(a.Etiketten.Any(), a.Parent.Etiketten.Any()) {
		errors.Log().Print("parent is not prefix joint")
		return
	}

	aEtt := a.Etiketten.Any()
	pEtt := a.Parent.Etiketten.Any()

	if aEtt.Equals(pEtt) {
		errors.Log().Print("parent is is equal to child")
		return
	}

	var ls kennung.Etikett

	if ls, err = kennung.LeftSubtract(aEtt, pEtt); err != nil {
		err = errors.Wrap(err)
		return
	}

	a.Etiketten = kennung.MakeEtikettSet(ls)

	return
}

// passed-in assignment may be nil?
func (atc *Refiner) Refine(a *Assignment) (err error) {
	if !atc.Enabled {
		return
	}

	if a.IsRoot {
		for _, c := range a.Children {
			if err = atc.Refine(c); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}

	if atc.shouldMergeIntoParent(a) {
		errors.Log().Print("merging into parent")
		p := a.Parent

		if err = p.consume(a); err != nil {
			err = errors.Wrap(err)
			return
		}

		return atc.Refine(p)
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

	if a.Etiketten == nil || a.Etiketten.Len() == 0 {
		return
	}

	childPrefixes := atc.childPrefixes(a)

	if len(childPrefixes) == 0 {
		return
	}

	groupingPrefix := childPrefixes[0]

	var na *Assignment

	if a.Etiketten.Len() == 1 &&
		a.Etiketten.Any().Equals(groupingPrefix.Etikett) {
		na = a
	} else {
		na = newAssignment(a.GetDepth() + 1)
		na.Etiketten = kennung.MakeEtikettSet(groupingPrefix.Etikett)
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

		c.Etiketten = kennung.SubtractPrefix(
			c.Etiketten,
			groupingPrefix.Etikett,
		)
	}

	return
}

type etikettBag struct {
	kennung.Etikett
	assignments []*Assignment
}

func (a Refiner) childPrefixes(node *Assignment) (out []etikettBag) {
	m := make(map[string][]*Assignment)
	out = make([]etikettBag, 0, len(node.Children))

	if node.Etiketten.Len() == 0 {
		return
	}

	for _, c := range node.Children {
		expanded := kennung.Expanded(c.Etiketten, expansion.ExpanderRight)

		expanded.Each(
			func(e kennung.Etikett) (err error) {
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
			var e1 kennung.Etikett

			errors.PanicIfError(e1.Set(e))

			out = append(out, etikettBag{Etikett: e1, assignments: n})
		}
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			if len(out[i].assignments) == len(out[j].assignments) {
				return len(
					out[i].Etikett.String(),
				) > len(
					out[j].Etikett.String(),
				)
			} else {
				return len(out[i].assignments) > len(out[j].assignments)
			}
		},
	)

	return
}
