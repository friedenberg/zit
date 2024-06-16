package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type PrefixSet struct {
	count    int
	innerMap map[string]objSet
}

type Segments struct {
	Ungrouped objSet
	Grouped   PrefixSet
}

func MakePrefixSet(c int) (s PrefixSet) {
	s.innerMap = make(map[string]objSet, c)
	return s
}

func MakePrefixSetFrom(
	ts objSet,
) (s PrefixSet) {
	s = MakePrefixSet(ts.Len())
	ts.Each(s.Add)
	return
}

func (s PrefixSet) Len() int {
	return s.count
}

func (s *PrefixSet) AddTransacted(z *sku.Transacted) (err error) {
	var o obj

	if err = o.Transacted.SetFromSkuLike(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.Add(&o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// this splits on right-expanded
func (s *PrefixSet) Add(z *obj) (err error) {
	es := kennung.Expanded(
		z.GetMetadatei().Verzeichnisse.GetImplicitEtiketten(),
		expansion.ExpanderRight,
	).CloneMutableSetPtrLike()

	if err = z.GetMetadatei().Verzeichnisse.GetExpandedEtiketten().EachPtr(
		es.AddPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if es.Len() == 0 {
		s.addPair("", z)
		return
	}

	if err = es.Each(
		func(e kennung.Etikett) (err error) {
			s.addPair(e.String(), z)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a PrefixSet) Subtract(
	b objSet,
) (c PrefixSet) {
	c = MakePrefixSet(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z *obj) (err error) {
				if b.Contains(z) {
					return
				}

				c.addPair(e, z)

				return
			},
		)
	}

	return
}

func (s *PrefixSet) addPair(
	e string,
	z *obj,
) {
	if e == z.Kennung.String() {
		e = ""
	}
	// if e != "" {
	// 	errors.PanicIfError((&kennung.Etikett{}).Set(e))
	// }

	existingSet, ok := s.innerMap[e]

	if !ok {
		existingSet = makeObjSet()
		s.innerMap[e] = existingSet
	}

	var existingObj *obj
	existingObj, ok = existingSet.Get(existingSet.Key(z))

	if ok && existingObj.IsDirectOrSelf() {
		z.SetDirect()
	} else if !ok {
		s.count += 1
	}

	existingSet.Add(z)
}

func (a PrefixSet) Each(
	f func(kennung.Etikett, objSet) error,
) (err error) {
	for e, ssz := range a.innerMap {
		var e1 kennung.Etikett

		if e != "" {
			e1 = kennung.MustEtikett(e)
		}

		if err = f(e1, ssz); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a PrefixSet) EachZettel(
	f schnittstellen.FuncIter[*obj],
) error {
	return a.Each(
		func(_ kennung.Etikett, st objSet) (err error) {
			err = st.Each(
				func(z *obj) (err error) {
					err = f(z)
					return
				},
			)

			return
		},
	)
}

func (a PrefixSet) Match(
	e kennung.Etikett,
) (out Segments) {
	out.Ungrouped = makeObjSet()
	out.Grouped = MakePrefixSet(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1 == "" {
			continue
		}

		zSet.Each(
			func(z *obj) (err error) {
				es := z.GetEtiketten()
				// var es kennung.EtikettSet

				// if es, err = allEtiketten(z); err != nil {
				// 	err = errors.Wrap(err)
				// 	return
				// }

				intersection := kennung.IntersectPrefixes(
					es,
					e,
				)

				exactMatch := intersection.Len() == 1 &&
					intersection.Any().Equals(e)

				if intersection.Len() == 0 && !exactMatch {
					return
				}

				for _, e2 := range iter.Elements(intersection) {
					out.Grouped.addPair(e2.String(), z)
				}

				return
			},
		)
	}

	return
}

func (a PrefixSet) Subset(
	e kennung.Etikett,
) (out Segments) {
	out.Ungrouped = makeObjSet()
	out.Grouped = MakePrefixSet(len(a.innerMap))

	e2 := catgut.MakeFromString(e.String())

	for e1, zSet := range a.innerMap {
		if e1 == "" {
			continue
		}

		zSet.Each(
			func(z *obj) (err error) {
				intersection := z.Metadatei.Verzeichnisse.Etiketten.All.GetMatching(e2)
				exactMatch := len(intersection) == 1 && intersection[0].Equals(e2)

				if len(intersection) > 0 && !exactMatch {
					for _, e2Match := range intersection {
						for _, e3 := range e2Match.Parents {
							out.Grouped.addPair(e3.Last().String(), z.cloneWithType(e3.Type))
						}
					}
				} else {
					out.Ungrouped.Add(z)
				}

				return
			},
		)
	}

	return
}
