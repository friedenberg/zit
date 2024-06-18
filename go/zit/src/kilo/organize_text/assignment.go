package organize_text

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

func newAssignment(d int) *Assignment {
	return &Assignment{
		Depth:     d,
		Etiketten: kennung.MakeEtikettSet(),
		objekten:  make(map[string]struct{}),
		Objekten:  make(Objekten, 0),
		Children:  make([]*Assignment, 0),
	}
}

type Assignment struct {
	IsRoot    bool
	Depth     int
	Etiketten kennung.EtikettSet
	objekten  map[string]struct{}
	Objekten
	Children []*Assignment
	Parent   *Assignment
}

func (a *Assignment) AddObjekte(v *obj) (err error) {
	k := key(&v.Transacted)
	_, ok := a.objekten[k]

	if ok {
		return
	}

	a.objekten[k] = struct{}{}

	return a.Objekten.Add(v)
}

func (a Assignment) GetDepth() int {
	if a.Parent == nil {
		return 0
	} else {
		return a.Parent.GetDepth() + 1
	}
}

func (a Assignment) MaxDepth() (d int) {
	d = a.GetDepth()

	for _, c := range a.Children {
		cd := c.MaxDepth()

		if d < cd {
			d = cd
		}
	}

	return
}

func (a Assignment) AlignmentSpacing() int {
	if a.Etiketten.Len() == 1 && kennung.IsDependentLeaf(a.Etiketten.Any()) {
		return a.Parent.AlignmentSpacing() + len(
			a.Parent.Etiketten.Any().String(),
		)
	}

	return 0
}

func (a Assignment) MaxLen() (m int) {
	a.Objekten.Each(
		func(z *obj) (err error) {
			oM := z.Kennung.Len()

			if oM > m {
				m = oM
			}

			return
		},
	)

	for _, c := range a.Children {
		oM := c.MaxLen()

		if oM > m {
			m = oM
		}
	}

	return
}

func (a Assignment) MaxKopfUndSchwanz(
	o Options,
) (kopf, schwanz int) {
	a.Objekten.Each(
		func(z *obj) (err error) {
			oKopf, oSchwanz := z.Kennung.LenKopfUndSchwanz()

			if o.PrintOptions.Abbreviations.Hinweisen {
				if oKopf, oSchwanz, err = o.Abbr.LenKopfUndSchwanz(
					&z.Kennung,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if oKopf > kopf {
				kopf = oKopf
			}

			if oSchwanz > schwanz {
				schwanz = oSchwanz
			}

			return
		},
	)

	for _, c := range a.Children {
		zKopf, zSchwanz := c.MaxKopfUndSchwanz(o)

		if zKopf > kopf {
			kopf = zKopf
		}

		if zSchwanz > schwanz {
			schwanz = zSchwanz
		}
	}

	return
}

func (a Assignment) String() (s string) {
	if a.Parent != nil {
		s = a.Parent.String() + "."
	}

	return s + iter.StringCommaSeparated(a.Etiketten)
}

func (a *Assignment) makeChild(e kennung.Etikett) (b *Assignment) {
	b = newAssignment(a.GetDepth() + 1)
	b.Etiketten = kennung.MakeEtikettSet(e)
	a.addChild(b)
	return
}

func (a *Assignment) makeChildWithSet(es kennung.EtikettSet) (b *Assignment) {
	b = newAssignment(a.GetDepth() + 1)
	b.Etiketten = es
	a.addChild(b)
	return
}

func (a *Assignment) addChild(c *Assignment) {
	if a == c {
		panic("child and parent are the same")
	}

	if c.Parent != nil && c.Parent == a {
		panic("child already has self as parent")
	}

	if c.Parent != nil {
		panic("child already has a parent")
	}

	a.Children = append(a.Children, c)
	c.Parent = a
}

func (a *Assignment) parentOrRoot() (p *Assignment) {
	switch a.Parent {
	case nil:
		return a

	default:
		return a.Parent
	}
}

func (a *Assignment) nthParent(n int) (p *Assignment, err error) {
	if n < 0 {
		n = -n
	}

	if n == 0 {
		p = a
		return
	}

	if a.Parent == nil {
		err = errors.Errorf("cannot get nth parent as parent is nil")
		return
	}

	return a.Parent.nthParent(n - 1)
}

func (a *Assignment) removeFromParent() (err error) {
	return a.Parent.removeChild(a)
}

func (a *Assignment) removeChild(c *Assignment) (err error) {
	if c.Parent != a {
		err = errors.Errorf("attempting to remove child from wrong parent")
		return
	}

	if len(a.Children) == 0 {
		err = errors.Errorf(
			"attempting to remove child when there are no children",
		)
		return
	}

	cap1 := 0
	cap2 := len(a.Children) - 1

	if cap2 > 0 {
		cap1 = cap2
	}

	nc := make([]*Assignment, 0, cap1)

	for _, c1 := range a.Children {
		if c1 == c {
			continue
		}

		nc = append(nc, c1)
	}

	c.Parent = nil
	a.Children = nc

	return
}

func (a *Assignment) consume(b *Assignment) (err error) {
	for _, c := range b.Children {
		if err = c.removeFromParent(); err != nil {
			err = errors.Wrap(err)
			return
		}

		a.addChild(c)
	}

	b.Objekten.Each(a.AddObjekte)
	b.Objekten.Each(b.Objekten.Del)

	if err = b.removeFromParent(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Assignment) AllEtiketten(mes kennung.EtikettMutableSet) (err error) {
	if a == nil {
		return
	}

	var es kennung.EtikettSet

	if es, err = a.expandedEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = es.EachPtr(mes.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Parent.AllEtiketten(mes); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Assignment) expandedEtiketten() (es kennung.EtikettSet, err error) {
	es = kennung.MakeEtikettSet()

	if a.Etiketten == nil {
		panic("etiketten are nil")
	}

	if a.Etiketten.Len() != 1 || a.Parent == nil {
		es = a.Etiketten.CloneSetPtrLike()
		return
	} else {
		e := a.Etiketten.Any()

		if kennung.IsDependentLeaf(e) {
			var pe kennung.EtikettSet

			if pe, err = a.Parent.expandedEtiketten(); err != nil {
				err = errors.Wrap(err)
				return
			}

			if pe.Len() > 1 {
				err = errors.Errorf(
					"cannot infer full etikett for assignment because parent assignment has more than one etiketten: %s",
					a.Parent.Etiketten,
				)

				return
			}

			e1 := pe.Any()

			if kennung.IsEmpty(e1) {
				err = errors.Errorf("parent etikett is empty")
				return
			}

			if err = e.Set(fmt.Sprintf("%s%s", e1, e)); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		es = kennung.MakeEtikettSet(e)
	}

	return
}

func (a *Assignment) SubtractFromSet(es kennung.EtikettMutableSet) (err error) {
	if err = a.Etiketten.EachPtr(
		func(e *kennung.Etikett) (err error) {
			if err = es.EachPtr(
				func(e1 *kennung.Etikett) (err error) {
					if !kennung.Contains(e1, e) {
						return
					}

					return es.DelPtr(e1)
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return es.DelPtr(e)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if a.Parent == nil {
		return
	}

	return a.Parent.SubtractFromSet(es)
}

func (a *Assignment) Contains(e *kennung.Etikett) bool {
	if a.Etiketten.ContainsKey(e.String()) {
		return true
	}

	if a.Parent == nil {
		return false
	}

	return a.Parent.Contains(e)
}

func (parent *Assignment) SortChildren() {
	sort.Slice(parent.Children, func(i, j int) bool {
		esi := parent.Children[i].Etiketten
		esj := parent.Children[j].Etiketten

		if esi.Len() == 1 && esj.Len() == 1 {
			ei := strings.TrimPrefix(esi.Any().String(), "-")
			ej := strings.TrimPrefix(esj.Any().String(), "-")

			ii, ierr := strconv.ParseInt(ei, 0, 64)
			ij, jerr := strconv.ParseInt(ej, 0, 64)

			if ierr == nil && jerr == nil {
				return ii < ij
			} else {
				return ei < ej
			}
		} else {
			vi := iter.StringCommaSeparated(esi)
			vj := iter.StringCommaSeparated(esj)
			return vi < vj
		}
	})
}
