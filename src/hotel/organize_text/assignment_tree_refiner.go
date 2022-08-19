package organize_text

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type AssignmentTreeRefiner struct {
	UsePrefixJoints bool
}

func (atc *AssignmentTreeRefiner) Refine(a *assignment) (err error) {
	logz.Print(a)
	sort.Slice(a.children, func(i, j int) bool {
		return a.children[i].etiketten.String() < a.children[j].etiketten.String()
	})

	if atc.UsePrefixJoints {
		if err = atc.applyPrefixJoints(a); err != nil {
			err = errors.Error(err)
			return
		}
	}

	for _, child := range a.children {
		// if i > 0 {
		// 	if child.etiketten.String() == a.children[i-1].etiketten.String() {
		// 		sib := a.children[i-1]

		// 		for _, c := range child.children {
		// 			c.parent = nil
		// 			sib.addChild(c)
		// 		}

		// 		child.parent = nil
		// 		child = sib

		// 		continue
		// 	}
		// }

		if err = atc.Refine(child); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

type etikettBag struct {
	etikett.Etikett
	assignments []*assignment
}

func (atc AssignmentTreeRefiner) applyPrefixJoints(a *assignment) (err error) {
	childPrefixes := atc.childPrefixes(a)

	if len(childPrefixes) > 0 {
		groupingPrefix := childPrefixes[0]

		na := newAssignment(a.depth + 1)
		na.etiketten = etikett.MakeSet(groupingPrefix.Etikett)
		a.addChild(na)

		for _, c := range groupingPrefix.assignments {
			if err = c.removeFromParent(); err != nil {
				err = errors.Error(err)
				return
			}

			c.etiketten = c.etiketten.SubtractPrefix(groupingPrefix.Etikett)
			c.depth += 1
			na.addChild(c)
		}
	}

	return
}

func (a AssignmentTreeRefiner) childPrefixes(node *assignment) (out []etikettBag) {
	logz.Print(node)

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
