package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/kennung"
)

type Refiner struct {
	Enabled         bool
	UsePrefixJoints bool
}

func (atc *Refiner) shouldMergeAllChildrenIntoParent(a *assignment) (ok bool) {
	switch {
	case a.parent.isRoot:
		fallthrough

	default:
		ok = false
	}

	return
}

func (atc *Refiner) shouldMergeIntoParent(a *assignment) bool {
	errors.Log().Printf("checking node should merge: %s", a)

	if a.parent == nil {
		errors.Log().Print("parent is nil")
		return false
	}

	if a.parent.isRoot {
		errors.Log().Print("parent is root")
		return false
	}

	if a.etiketten.Len() == 1 && a.etiketten.Any().IsEmpty() {
		errors.Log().Print("1 Etikett, and it's empty, merging")
		return true
	}

	if a.etiketten.Len() == 0 {
		errors.Log().Print("etiketten length is 0, merging")
		return true
	}

	if a.parent.etiketten.Len() != 1 {
		errors.Log().Print("parent etiketten length is not 1")
		return false
	}

	if a.etiketten.Len() != 1 {
		errors.Log().Print("etiketten length is not 1")
		return false
	}

	if !a.etiketten.Equals(a.parent.etiketten) {
		errors.Log().Print("parent etiketten not equal")
		return false
	}

	if kennung.IsDependentLeaf(a.parent.etiketten.Any()) {
		errors.Log().Print("is prefix joint")
		return false
	}

	if kennung.IsDependentLeaf(a.etiketten.Any()) {
		errors.Log().Print("is prefix joint")
		return false
	}

	return true
}

func (atc *Refiner) renameForPrefixJoint(a *assignment) (err error) {
	if !atc.UsePrefixJoints {
		return
	}

	if a == nil {
		errors.Log().Printf("assignment is nil")
		return
	}

	if a.parent == nil {
		errors.Log().Printf("parent is nil: %#v", a)
		return
	}

	if a.parent.etiketten.Len() == 0 {
		return
	}

	if a.parent.etiketten.Len() != 1 {
		return
	}

	if kennung.IsDependentLeaf(a.parent.etiketten.Any()) {
		return
	}

	if kennung.IsDependentLeaf(a.etiketten.Any()) {
		return
	}

	if !kennung.HasParentPrefix(a.etiketten.Any(), a.parent.etiketten.Any()) {
		errors.Log().Print("parent is not prefix joint")
		return
	}

	var ls kennung.Etikett

	if ls, err = a.etiketten.Any().LeftSubtract(a.parent.etiketten.Any()); err != nil {
		err = errors.Wrap(err)
		return
	}

	a.etiketten = kennung.MakeEtikettSet(ls)

	return
}

// passed-in assignment may be nil?
func (atc *Refiner) Refine(a *assignment) (err error) {
	if !atc.Enabled {
		return
	}

	if a.isRoot {
		for _, c := range a.children {
			if err = atc.Refine(c); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}

	//TODO fix after breaking during migration to collections
	// if atc.shouldMergeIntoParent(a) {
	// 	errors.Log().Print("merging into parent")
	// 	p := a.parent

	// 	if err = p.consume(a); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	return atc.Refine(p)
	// }

	// if err = atc.applyPrefixJoints(a); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if err = atc.renameForPrefixJoint(a); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, child := range a.children {
		if err = atc.Refine(child); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// if err = atc.applyPrefixJoints(a); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	sort.Slice(a.children, func(i, j int) bool {
		return a.children[i].etiketten.String() < a.children[j].etiketten.String()
	})

	return
}

func (atc Refiner) applyPrefixJoints(a *assignment) (err error) {
	if !atc.UsePrefixJoints {
		return
	}

	if a.etiketten.Len() == 0 {
		return
	}

	childPrefixes := atc.childPrefixes(a)

	if len(childPrefixes) == 0 {
		return
	}

	groupingPrefix := childPrefixes[0]

	var na *assignment

	if a.etiketten.Len() == 1 && a.etiketten.Any().Equals(&groupingPrefix.Etikett) {
		na = a
	} else {
		na = newAssignment(a.Depth() + 1)
		na.etiketten = kennung.MakeEtikettSet(groupingPrefix.Etikett)
		a.addChild(na)
	}

	for _, c := range groupingPrefix.assignments {
		if c.parent != na {
			if err = c.removeFromParent(); err != nil {
				err = errors.Wrap(err)
				return
			}

			na.addChild(c)
		}

		c.etiketten = kennung.SubtractPrefix(c.etiketten, groupingPrefix.Etikett)
	}

	return
}

type etikettBag struct {
	kennung.Etikett
	assignments []*assignment
}

func (a Refiner) childPrefixes(node *assignment) (out []etikettBag) {
	m := make(map[kennung.Etikett][]*assignment)
	out = make([]etikettBag, 0, len(node.children))

	if node.etiketten.Len() == 0 {
		return
	}

	for _, c := range node.children {
		expanded := kennung.Expanded(c.etiketten, kennung.ExpanderRight)

		expanded.Each(
			func(e kennung.Etikett) (err error) {
				if e.String() == "" {
					return
				}

				var n []*assignment
				ok := false

				if n, ok = m[e]; !ok {
					n = make([]*assignment, 0)
				}

				n = append(n, c)

				m[e] = n

				return
			},
		)
	}

	for e, n := range m {
		if len(n) > 1 {
			out = append(out, etikettBag{Etikett: e, assignments: n})
		}
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			if len(out[i].assignments) == len(out[j].assignments) {
				return len(out[i].Etikett.String()) > len(out[j].Etikett.String())
			} else {
				return len(out[i].assignments) > len(out[j].assignments)
			}
		},
	)

	return
}
