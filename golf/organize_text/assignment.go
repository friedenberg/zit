package organize_text

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/log"
	"github.com/friedenberg/zit/charlie/etikett"
)

type assignment struct {
	etiketten etikett.Set
	named     zettelSet
	unnamed   newZettelSet
	depth     int
	children  []*assignment
	parent    *assignment
}

func newAssignment(depth int) *assignment {
	return &assignment{
		etiketten: etikett.MakeSet(),
		named:     makeZettelSet(),
		unnamed:   makeNewZettelSet(),
		depth:     depth,
		children:  make([]*assignment, 0),
	}
}

func (a assignment) String() (s string) {
	if a.parent != nil {
		log.PrintDebug(a)
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
