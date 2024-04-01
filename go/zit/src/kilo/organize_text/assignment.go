package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/compare_map"
)

type Assignment struct {
	IsRoot    bool
	Depth     int
	Etiketten kennung.EtikettSet
	Named     schnittstellen.MutableSetLike[*obj]
	Unnamed   schnittstellen.MutableSetLike[*obj]
	Children  []*Assignment
	Parent    *Assignment
}

func newAssignment(d int) *Assignment {
	return &Assignment{
		Depth:     d,
		Etiketten: kennung.MakeEtikettSet(),
		Named:     collections_value.MakeMutableValueSet[*obj](nil),
		Unnamed:   collections_value.MakeMutableValueSet[*obj](nil),
		Children:  make([]*Assignment, 0),
	}
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
	a.Named.Each(
		func(z *obj) (err error) {
			oM := z.Sku.Kennung.Len()

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

func (a Assignment) MaxKopfUndSchwanz() (kopf, schwanz int) {
	a.Named.Each(
		func(z *obj) (err error) {
			oKopf, oSchwanz := z.Sku.Kennung.LenKopfUndSchwanz()

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

func (a Assignment) String() (s string) {
	if a.Parent != nil {
		s = a.Parent.String() + "."
	}

	return s + iter.StringCommaSeparated[kennung.Etikett](a.Etiketten)
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

	b.Named.Each(a.Named.Add)
	b.Named.Each(b.Named.Del)

	b.Unnamed.Each(a.Unnamed.Add)
	b.Unnamed.Each(b.Unnamed.Del)

	if err = b.removeFromParent(); err != nil {
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

func (a *Assignment) addToCompareMap(
	ot *Text,
	m Metadatei,
	es kennung.EtikettSet,
	out *compare_map.CompareMap,
) (err error) {
	mes := es.CloneMutableSetPtrLike()

	var es1 kennung.EtikettSet

	if es1, err = a.expandedEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	es1.Each(mes.Add)
	es = mes.CloneSetPtrLike()

	if err = a.Named.Each(
		func(z *obj) (err error) {
			if z.Sku.Kennung.String() == "" {
				panic(fmt.Sprintf("%s: Kennung is nil", z))
			}

			fk := kennung.FormattedString(&z.Sku.Kennung)
			out.Named.Add(fk, z.Sku.Metadatei.Bezeichnung)

			for _, e := range iter.SortedValues[kennung.Etikett](es) {
				out.Named.AddEtikett(fk, e, z.Sku.Metadatei.Bezeichnung)
			}

			for _, e := range iter.Elements[kennung.Etikett](m.EtikettSet) {
				errors.TodoP4("add typ")
				out.Named.AddEtikett(fk, e, z.Sku.Metadatei.Bezeichnung)
			}

			if ot.Konfig.NewOrganize {
				if err = z.Sku.Metadatei.GetEtiketten().EachPtr(
					func(e *kennung.Etikett) (err error) {
						if a.Contains(e) {
							return
						}

						out.Named.AddEtikett(fk, *e, z.Sku.Metadatei.Bezeichnung)
						return
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Unnamed.Each(
		func(z *obj) (err error) {
			out.Unnamed.Add(
				z.Sku.Metadatei.Bezeichnung.String(),
				z.Sku.Metadatei.Bezeichnung,
			)

			for _, e := range iter.SortedValues[kennung.Etikett](es) {
				out.Unnamed.AddEtikett(
					z.Sku.Metadatei.Bezeichnung.String(),
					e,
					z.Sku.Metadatei.Bezeichnung,
				)
			}

			for _, e := range iter.Elements[kennung.Etikett](m.EtikettSet) {
				errors.TodoP4("add typ")
				out.Unnamed.AddEtikett(
					z.Sku.Metadatei.Bezeichnung.String(),
					e,
					z.Sku.Metadatei.Bezeichnung,
				)
			}

			if ot.Konfig.NewOrganize {
				if err = z.Sku.Metadatei.GetEtiketten().EachPtr(
					func(e *kennung.Etikett) (err error) {
						if a.Contains(e) {
							return
						}

						out.Named.AddEtikett(
							z.Sku.Metadatei.Bezeichnung.String(),
							*e,
							z.Sku.Metadatei.Bezeichnung,
						)

						return
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, c := range a.Children {
		if err = c.addToCompareMap(ot, m, es, out); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
