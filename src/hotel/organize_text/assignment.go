package organize_text

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type assignment struct {
	etiketten etikett.Set
	named     zettelSet
	unnamed   newZettelSet
	children  []*assignment
	parent    *assignment
}

func newAssignment() *assignment {
	return &assignment{
		etiketten: etikett.MakeSet(),
		named:     makeZettelSet(),
		unnamed:   makeNewZettelSet(),
		children:  make([]*assignment, 0),
	}
}

func (a assignment) Depth() int {
	if a.parent == nil {
		return 0
	} else {
		return a.parent.Depth() + 1
	}
}

func (a assignment) String() (s string) {
	if a.parent != nil {
		s = a.parent.String() + "."
	}

	return s + a.etiketten.String()
}

func (a *assignment) addChild(c *assignment) {
	if a == c {
		panic("child and parent are the same")
	}

	if c.parent == a {
		panic("child already has self as parent")
	}

	if c.parent != nil {
		panic("child already has a parent")
	}

	a.children = append(a.children, c)
	c.parent = a
}

func (a *assignment) nthParent(n int) (p *assignment, err error) {
	if n < 0 {
		n = -n
	}

	if n == 0 {
		p = a
		return
	}

	if a.parent == nil {
		err = errors.Errorf("cannot get nth parent as parent is nil")
		return
	}

	return a.parent.nthParent(n - 1)
}

func (a *assignment) removeFromParent() (err error) {
	return a.parent.removeChild(a)
}

func (a *assignment) removeChild(c *assignment) (err error) {
	if c.parent != a {
		err = errors.Errorf("attempting to remove child from wrong parent")
		return
	}

	if len(a.children) == 0 {
		err = errors.Errorf("attempting to remove child when there are no children")
		return
	}

	cap1 := 0
	cap2 := len(c.children) - 1

	if cap2 > 0 {
		cap1 = cap2
	}

	nc := make([]*assignment, 0, cap1)

	for _, c1 := range a.children {
		if c1 == c {
			continue
		}

		nc = append(nc, c1)
	}

	c.parent = nil
	a.children = nc

	return
}

func (a *assignment) consume(b *assignment) (err error) {
	for _, c := range b.children {
		if err = c.removeFromParent(); err != nil {
			err = errors.Error(err)
			return
		}

		a.parent.addChild(c)
	}

	for k, v := range b.named {
		a.parent.named[k] = v
	}

	for k, v := range b.unnamed {
		a.parent.unnamed[k] = v
	}

	if err = b.removeFromParent(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (a *assignment) expandedEtiketten() (es etikett.Set, err error) {
	logz.Print(a.etiketten)
	if a.etiketten.Len() != 1 || a.parent == nil {
		es = *(a.etiketten.Copy())
		return
	} else {
		e := a.etiketten.Any()

		if e.IsDependentLeaf() {
			var pe etikett.Set

			if pe, err = a.parent.expandedEtiketten(); err != nil {
				err = errors.Error(err)
				return
			}

			if pe.Len() > 1 {
				err = errors.Errorf(
					"cannot infer full etikett for assignment because parent assignment has more than one etiketten: %s",
					a.parent.etiketten,
				)

				return
			}

			e1 := pe.Any()

			if e1.IsEmpty() {
				err = errors.Errorf("parent etikett is empty")
				return
			}

			logz.Print(e1, e)
			e = etikett.Etikett{Value: fmt.Sprintf("%s%s", e1, e)}
			logz.Print(e)
		}

		es = etikett.MakeSet(e)
		logz.Print(es)
	}

	return
}
