package zettel_named

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/kennung"
)

type SetPrefixNamed map[kennung.Etikett]MutableSet

type SetPrefixNamedSegments struct {
	Ungrouped MutableSet
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	s := make(SetPrefixNamed)
	return &s
}

// this splits on right-expanded
func (s *SetPrefixNamed) Add(z Zettel) {
	es := kennung.Expanded(z.Stored.Objekte.Etiketten, kennung.ExpanderRight)

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (s *SetPrefixNamed) addPair(e kennung.Etikett, z Zettel) {
	existing, ok := (*s)[e]

	if !ok {
		existing = MakeMutableSet()
	}

	existing.Add(&z)
	(*s)[e] = existing
}

// for all of the zettels, check for intersections with the passed in
// etikett, and if there is a prefix match, group it out the output set segments
// appropriately
func (a SetPrefixNamed) Subset(e kennung.Etikett) (out SetPrefixNamedSegments) {
	out.Ungrouped = MakeMutableSet()
	out.Grouped = NewSetPrefixNamed()

	for e1, zSet := range a {
		zSet.Each(
			func(z *Zettel) (err error) {
				intersection := z.Stored.Objekte.Etiketten.IntersectPrefixes(kennung.MakeEtikettSet(e))
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

func (s SetPrefixNamed) ToSetNamed() (out MutableSet) {
	out = MakeMutableSet()

	for _, zs := range s {
		zs.Each(
			func(z *Zettel) (err error) {
				out.Add(z)

				return
			},
		)
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
