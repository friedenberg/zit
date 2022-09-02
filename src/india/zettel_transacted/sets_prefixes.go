package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type SetPrefixTransacted struct {
	innerMap map[etikett.Etikett]Set
}

type SetPrefixTransactedSegments struct {
	Ungrouped Set
	Grouped   SetPrefixTransacted
}

func MakeSetPrefixTransacted(c int) (s SetPrefixTransacted) {
	s.innerMap = make(map[etikett.Etikett]Set, c)
	return s
}

//TODO mark that this splits on right-expanded
func (s *SetPrefixTransacted) Add(z Zettel) {
	es := z.Named.Stored.Zettel.Etiketten.Expanded(etikett.ExpanderRight{})

	if es.Len() == 0 {
		es = etikett.MakeSet(etikett.Etikett{})
	}

	for _, e := range es {
		s.addPair(e, z)
	}
}

func (s SetPrefixTransacted) addPair(e etikett.Etikett, z Zettel) {
	existing, ok := s.innerMap[e]

	if !ok {
		existing = MakeSetUnique(1)
	}

	existing.Add(z)
	s.innerMap[e] = existing
}

func (a SetPrefixTransacted) Each(f func(etikett.Etikett, Set) error) (err error) {
	for e, ssz := range a.innerMap {
		if err = f(e, ssz); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			} else {
				err = errors.Error(err)
			}

			return
		}
	}

	return
}

func (a SetPrefixTransacted) EachZettel(f func(etikett.Etikett, Zettel) error) (err error) {
	return a.Each(
		func(e etikett.Etikett, st Set) (err error) {
			for _, sz := range st.innerMap {
				if err = f(e, sz); err != nil {
					return
				}
			}

			return
		},
	)
}

// for all of the zettels, check for intersections with the passed in
// etikett, and if there is a prefix match, group it out the output set segments
// appropriately
func (a SetPrefixTransacted) Subset(e etikett.Etikett) (out SetPrefixTransactedSegments) {
	out.Ungrouped = MakeSetUnique(len(a.innerMap))
	out.Grouped = MakeSetPrefixTransacted(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		for _, z := range zSet.innerMap {
			intersection := z.Named.Stored.Zettel.Etiketten.IntersectPrefixes(etikett.MakeSet(e))
			logz.Printf("%s yields %s", e1, intersection)

			if intersection.Len() > 0 {
				for _, e2 := range intersection {
					out.Grouped.addPair(e2, z)
				}
			} else {
				out.Ungrouped.Add(z)
			}
		}
	}

	return
}

func (s SetPrefixTransacted) ToSet() (out Set) {
	out = MakeSetUnique(len(s.innerMap))

	for _, zs := range s.innerMap {
		for _, z := range zs.innerMap {
			out.Add(z)
		}
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
