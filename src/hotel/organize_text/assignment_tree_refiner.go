package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type AssignmentTreeRefiner struct {
	Enabled         bool
	UsePrefixJoints bool
}

func (atc *AssignmentTreeRefiner) Refine(a *assignment) (err error) {
	if !atc.Enabled {
		return
	}

	if a.parent != nil && a.etiketten.Equals(a.parent.etiketten) {
		p := a.parent

		if err = p.consume(a); err != nil {
			err = errors.Error(err)
			return
		}

		return atc.Refine(p)
	}

	if atc.UsePrefixJoints {
		if err = atc.applyPrefixJoints(a); err != nil {
			err = errors.Error(err)
			return
		}
	}

	for _, child := range a.children {
		if err = atc.Refine(child); err != nil {
			err = errors.Error(err)
			return
		}
	}

	if atc.UsePrefixJoints {
		if err = atc.applyPrefixJoints(a); err != nil {
			err = errors.Error(err)
			return
		}
	}

	sort.Slice(a.children, func(i, j int) bool {
		return a.children[i].etiketten.String() < a.children[j].etiketten.String()
	})

	return
}

type etikettBag struct {
	etikett.Etikett
	assignments []*assignment
}

func (atc AssignmentTreeRefiner) applyPrefixJoints(a *assignment) (err error) {
	childPrefixes := atc.childPrefixes(a)

	if len(childPrefixes) == 0 {
		return
	}

	groupingPrefix := childPrefixes[0]

	var na *assignment

	if a.etiketten.Len() == 1 && a.etiketten.Any().Equals(groupingPrefix.Etikett) {
		na = a
	} else {
		na = newAssignment()
		na.etiketten = etikett.MakeSet(groupingPrefix.Etikett)
		a.addChild(na)
	}

	for _, c := range groupingPrefix.assignments {
		if c.parent != na {
			if err = c.removeFromParent(); err != nil {
				err = errors.Error(err)
				return
			}

			na.addChild(c)
		}

		c.etiketten = c.etiketten.SubtractPrefix(groupingPrefix.Etikett)
	}

	return
}

func (a AssignmentTreeRefiner) childPrefixes(node *assignment) (out []etikettBag) {
	m := make(map[etikett.Etikett][]*assignment)
	out = make([]etikettBag, 0, len(node.children))

	if node.etiketten.Len() == 0 {
		return
	}

	for _, c := range node.children {
		for _, e := range c.etiketten.Expanded(etikett.ExpanderRight{}) {
			if e.String() == "" {
				continue
			}

			var n []*assignment
			ok := false

			if n, ok = m[e]; !ok {
				n = make([]*assignment, 0)
			}

			n = append(n, c)

			m[e] = n
		}
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
