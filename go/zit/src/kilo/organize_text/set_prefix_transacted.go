package organize_text

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type PrefixSet struct {
	count    int
	innerMap map[string]sku.TransactedMutableSet
}

type Segments struct {
	Ungrouped sku.TransactedMutableSet
	Grouped   PrefixSet
}

func MakePrefixSet(c int) (s PrefixSet) {
	s.innerMap = make(map[string]sku.TransactedMutableSet, c)
	return s
}

func MakePrefixSetFrom(
	ts sku.TransactedSet,
) (s PrefixSet) {
	s = MakePrefixSet(ts.Len())
	ts.Each(s.Add)
	return
}

func (s PrefixSet) Len() int {
	return s.count
}

// this splits on right-expanded
func (s *PrefixSet) Add(z *sku.Transacted) (err error) {
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
		es = kennung.MakeEtikettMutableSet(kennung.Etikett{})
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
	b sku.TransactedMutableSet,
) (c PrefixSet) {
	c = MakePrefixSet(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z *sku.Transacted) (err error) {
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
	z *sku.Transacted,
) {
	s.count += 1

	existing, ok := s.innerMap[e]

	if !ok {
		existing = sku.MakeTransactedMutableSet()
	}

	existing.Add(z)
	s.innerMap[e] = existing
}

func (a PrefixSet) Each(
	f func(kennung.Etikett, sku.TransactedMutableSet) error,
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
	f schnittstellen.FuncIter[*sku.Transacted],
) error {
	return a.Each(
		func(_ kennung.Etikett, st sku.TransactedMutableSet) (err error) {
			err = st.Each(
				func(z *sku.Transacted) (err error) {
					err = f(z)
					return
				},
			)

			return
		},
	)
}

func (a PrefixSet) EachPair(
	f func(kennung.Etikett, *sku.Transacted) error,
) error {
	return a.Each(
		func(e kennung.Etikett, st sku.TransactedMutableSet) (err error) {
			err = st.Each(
				func(z *sku.Transacted) (err error) {
					err = f(e, z)
					return
				},
			)

			return
		},
	)
}

// for all of the zettels, check for intersections with the passed in
// etikett, and if there is a prefix match, group it out the output set segments
// appropriately

func (s PrefixSet) ToSet() (out sku.TransactedMutableSet) {
	out = sku.MakeTransactedMutableSet()

	for _, zs := range s.innerMap {
		zs.Each(out.Add)
	}

	return
}

func (a PrefixSet) Match(
	e kennung.Etikett,
) (out Segments) {
	out.Ungrouped = sku.MakeTransactedMutableSet()
	out.Grouped = MakePrefixSet(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1 == "" {
			continue
		}

		zSet.Each(
			func(z *sku.Transacted) (err error) {
				es := z.GetEtiketten()
				// var es kennung.EtikettSet

				// if es, err = allEtiketten(z); err != nil {
				// 	err = errors.Wrap(err)
				// 	return
				// }

				intersection := kennung.IntersectPrefixes(
					es,
					kennung.MakeEtikettSet(e),
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
	out.Ungrouped = sku.MakeTransactedMutableSet()
	out.Grouped = MakePrefixSet(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1 == "" {
			continue
		}

		zSet.Each(
			func(z *sku.Transacted) (err error) {
				es := z.GetEtiketten()
				// var es kennung.EtikettSet

				// if es, err = allEtiketten(z); err != nil {
				// 	err = errors.Wrap(err)
				// 	return
				// }

				intersection := kennung.IntersectPrefixes(
					es,
					kennung.MakeEtikettSet(e),
				)

				exactMatch := intersection.Len() == 1 &&
					intersection.Any().Equals(e)

				if intersection.Len() > 0 && !exactMatch {
					for _, e2 := range iter.Elements(intersection) {
						out.Grouped.addPair(e2.String(), z)
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
