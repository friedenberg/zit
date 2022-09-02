package organize_text

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type assignment struct {
	isRoot    bool
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

	if c.parent != nil && c.parent == a {
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
	errors.Print(a.etiketten)
	errors.Print(a.named)
	errors.Print(b.etiketten)
	errors.Print(b.named)
	errors.Caller(1, "test")

	for _, c := range b.children {
		errors.Print(c)
		if err = c.removeFromParent(); err != nil {
			err = errors.Error(err)
			return
		}

		errors.Print(a)
		a.addChild(c)
	}

	for k, v := range b.named {
		a.named[k] = v
		delete(b.named, k)
	}

	for k, v := range b.unnamed {
		a.unnamed[k] = v
		delete(b.unnamed, k)
	}

	if err = b.removeFromParent(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (a *assignment) expandedEtiketten() (es etikett.Set, err error) {
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

			errors.Print(e1, e)
			e = etikett.Etikett{Value: fmt.Sprintf("%s%s", e1, e)}
			errors.Print(e)
		}

		es = etikett.MakeSet(e)
		errors.Print(es)
	}

	return
}
