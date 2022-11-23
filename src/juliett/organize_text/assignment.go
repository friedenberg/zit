package organize_text

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/kennung"
)

type assignment struct {
	isRoot    bool
	depth     int
	etiketten kennung.Set
	named     collections.MutableValueSet[zettel, *zettel]
	unnamed   collections.MutableValueSet[newZettel, *newZettel]
	children  []*assignment
	parent    *assignment
}

func newAssignment(d int) *assignment {
	return &assignment{
		depth:     d,
		etiketten: kennung.MakeSet(),
		named:     collections.MakeMutableValueSet[zettel](),
		unnamed:   collections.MakeMutableValueSet[newZettel](),
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

func (a assignment) MaxDepth() (d int) {
	d = a.Depth()

	for _, c := range a.children {
		cd := c.MaxDepth()

		if d < cd {
			d = cd
		}
	}

	return
}

func (a assignment) AlignmentSpacing() int {
	if a.etiketten.Len() == 1 && a.etiketten.Any().IsDependentLeaf() {
		return a.parent.AlignmentSpacing() + len(a.parent.etiketten.Any().String())
	}

	return 0
}

func (a assignment) MaxKopfUndSchwanz() (kopf, schwanz int) {
	a.named.Each(
		func(z zettel) (err error) {
			parts := [2]string{z.Hinweis.Kopf(), z.Hinweis.Schwanz()}
			zKopf := len(parts[0])
			zSchwanz := len(parts[1])

			if zKopf > kopf {
				kopf = zKopf
			}

			if zSchwanz > schwanz {
				schwanz = zSchwanz
			}

			return
		},
	)

	for _, c := range a.children {
		zKopf, zSchwanz := c.MaxKopfUndSchwanz()

		if zKopf > kopf {
			kopf = zKopf
		}

		if zSchwanz > schwanz {
			schwanz = zSchwanz
		}
	}

	return
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

func (a *assignment) parentOrRoot() (p *assignment) {
	switch {
	case a.parent == nil:
		return a

	default:
		return a.parent
	}
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
	cap2 := len(a.children) - 1

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
			err = errors.Wrap(err)
			return
		}

		a.addChild(c)
	}

	b.named.Each(a.named.Add)
	b.named.Each(b.named.Del)

	b.unnamed.Each(a.unnamed.Add)
	b.unnamed.Each(b.unnamed.Del)

	if err = b.removeFromParent(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *assignment) expandedEtiketten() (es kennung.Set, err error) {
	es = kennung.MakeSet()

	if a.etiketten.Len() != 1 || a.parent == nil {
		es = a.etiketten.Copy()
		return
	} else {
		e := a.etiketten.Any()

		if e.IsDependentLeaf() {
			var pe kennung.Set

			if pe, err = a.parent.expandedEtiketten(); err != nil {
				err = errors.Wrap(err)
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
			e = kennung.Etikett{Value: fmt.Sprintf("%s%s", e1, e)}
			errors.Print(e)
		}

		es = kennung.MakeSet(e)
	}

	return
}
