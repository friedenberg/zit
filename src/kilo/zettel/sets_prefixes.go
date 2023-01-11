package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

type SetPrefixVerzeichnisse struct {
	count    int
	innerMap map[kennung.Etikett]MutableSet
}

type SetPrefixVerzeichnisseSegments struct {
	Ungrouped MutableSet
	Grouped   SetPrefixVerzeichnisse
}

func MakeSetPrefixVerzeichnisse(c int) (s SetPrefixVerzeichnisse) {
	s.innerMap = make(map[kennung.Etikett]MutableSet, c)
	return s
}

func (s SetPrefixVerzeichnisse) Len() int {
	return s.count
}

// this splits on right-expanded
func (s *SetPrefixVerzeichnisse) Add(z Transacted) {
	es := kennung.Expanded(z.Objekte.Etiketten, kennung.ExpanderRight)

	if es.Len() == 0 {
		es = kennung.MakeEtikettSet(kennung.Etikett{})
	}

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (a SetPrefixVerzeichnisse) Subtract(b MutableSet) (c SetPrefixVerzeichnisse) {
	c = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z *Transacted) (err error) {
				if b.Contains(z) {
					return
				}

				c.addPair(e, *z)

				return
			},
		)
	}

	return
}

func (s *SetPrefixVerzeichnisse) addPair(e kennung.Etikett, z Transacted) {
	s.count += 1

	existing, ok := s.innerMap[e]

	if !ok {
		existing = MakeMutableSetUnique(1)
	}

	existing.Add(&z)
	s.innerMap[e] = existing
}

func (a SetPrefixVerzeichnisse) Each(f func(kennung.Etikett, MutableSet) error) (err error) {
	for e, ssz := range a.innerMap {
		if err = f(e, ssz); err != nil {
			if errors.Is(err, collections.ErrStopIteration) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a SetPrefixVerzeichnisse) EachZettel(f func(kennung.Etikett, Transacted) error) error {
	return a.Each(
		func(e kennung.Etikett, st MutableSet) (err error) {
			st.Each(
				func(z *Transacted) (err error) {
					err = f(e, *z)
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
func (a SetPrefixVerzeichnisse) Subset(e kennung.Etikett) (out SetPrefixVerzeichnisseSegments) {
	out.Ungrouped = MakeMutableSetUnique(len(a.innerMap))
	out.Grouped = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1.String() == "" {
			continue
		}

		zSet.Each(
			func(z *Transacted) (err error) {
				intersection := z.Objekte.Etiketten.IntersectPrefixes(kennung.MakeEtikettSet(e))
				errors.Log().Printf("%s yields %s", e1, intersection)

				if intersection.Len() > 0 {
					for _, e2 := range intersection.Elements() {
						out.Grouped.addPair(e2, *z)
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

func (s SetPrefixVerzeichnisse) ToSet() (out MutableSet) {
	out = MakeMutableSetUnique(len(s.innerMap))

	for _, zs := range s.innerMap {
		zs.Each(out.Add)
	}

	return
}

// func (s SetNamed) Slice() (slice []string) {
// 	slice = make([]string, len(zs.etikettenToExisting))
// 	i := 0

// 	for e, _ := range zs.etikettenToExisting {
// 		sorted[i] = e
// 		i++
// 	}

// 	sort.Slice(sorted, func(i, j int) bool {
// 		return sorted[i] < sorted[j]
// 	})

// 	return
// }
