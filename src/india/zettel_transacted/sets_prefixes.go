package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type SetPrefixTransacted struct {
	count    int
	innerMap map[etikett.Etikett]MutableSet
}

type SetPrefixTransactedSegments struct {
	Ungrouped MutableSet
	Grouped   SetPrefixTransacted
}

func MakeSetPrefixTransacted(c int) (s SetPrefixTransacted) {
	s.innerMap = make(map[etikett.Etikett]MutableSet, c)
	return s
}

func (s SetPrefixTransacted) Len() int {
	return s.count
}

// this splits on right-expanded
func (s *SetPrefixTransacted) Add(z Zettel) {
	es := etikett.Expanded(z.Named.Stored.Zettel.Etiketten, etikett.ExpanderRight{})

	if es.Len() == 0 {
		es = etikett.MakeSet(etikett.Etikett{})
	}

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (a SetPrefixTransacted) Subtract(b MutableSet) (c SetPrefixTransacted) {
	c = MakeSetPrefixTransacted(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z *Zettel) (err error) {
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

func (s *SetPrefixTransacted) addPair(e etikett.Etikett, z Zettel) {
	s.count += 1

	existing, ok := s.innerMap[e]

	if !ok {
		existing = MakeMutableSetUnique(1)
	}

	existing.Add(&z)
	s.innerMap[e] = existing
}

func (a SetPrefixTransacted) Each(f func(etikett.Etikett, MutableSet) error) (err error) {
	for e, ssz := range a.innerMap {
		if err = f(e, ssz); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a SetPrefixTransacted) EachZettel(f func(etikett.Etikett, Zettel) error) error {
	return a.Each(
		func(e etikett.Etikett, st MutableSet) (err error) {
			st.Each(
				func(z *Zettel) (err error) {
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
func (a SetPrefixTransacted) Subset(e etikett.Etikett) (out SetPrefixTransactedSegments) {
	out.Ungrouped = MakeMutableSetUnique(len(a.innerMap))
	out.Grouped = MakeSetPrefixTransacted(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1.String() == "" {
			continue
		}

		zSet.Each(
			func(z *Zettel) (err error) {
				intersection := z.Named.Stored.Zettel.Etiketten.IntersectPrefixes(etikett.MakeSet(e))
				errors.Printf("%s yields %s", e1, intersection)

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

func (s SetPrefixTransacted) ToSet() (out MutableSet) {
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