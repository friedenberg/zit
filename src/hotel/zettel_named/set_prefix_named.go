package zettel_named

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
)

type SetPrefixNamed map[etikett.Etikett]SetNamed

type SetPrefixNamedSegments struct {
	Ungrouped *SetNamed
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	s := make(SetPrefixNamed)
	return &s
}

//TODO mark that this splits on right-expanded
func (s *SetPrefixNamed) Add(z Zettel) {
	es := z.Stored.Zettel.Etiketten.Expanded(etikett.ExpanderRight{})

	for _, e := range es {
		s.addPair(e, z)
	}
}

func (s *SetPrefixNamed) addPair(e etikett.Etikett, z Zettel) {
	existing, ok := (*s)[e]

	if !ok {
		existing = MakeSetNamed()
	}

	existing.Add(z)
	(*s)[e] = existing
}

// for all of the zettels, check for intersections with the passed in
// etikett, and if there is a prefix match, group it out the output set segments
// appropriately
func (a SetPrefixNamed) Subset(e etikett.Etikett) (out SetPrefixNamedSegments) {
	out.Ungrouped = NewSetNamed()
	out.Grouped = NewSetPrefixNamed()

	for e1, zSet := range a {
		for _, z := range zSet {
			intersection := z.Stored.Zettel.Etiketten.IntersectPrefixes(etikett.MakeSet(e))
			errors.Printf("%s yields %s", e1, intersection)

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

func (s SetPrefixNamed) ToSetNamed() (out *SetNamed) {
	out = NewSetNamed()

	for _, zs := range s {
		for _, z := range zs {
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
