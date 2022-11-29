package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/id_set"
)

//TODO rename
type SetPrefixNamed map[kennung.Etikett]collections.MutableSet[id_set.Element]

type SetPrefixNamedSegments struct {
	Ungrouped collections.MutableSet[id_set.Element]
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	s := make(SetPrefixNamed)
	return &s
}

func makeMutableZettelLikeSet() collections.MutableSet[id_set.Element] {
	return collections.MakeMutableSet(
		func(e id_set.Element) string {
			if e == nil {
				return ""
			}

			return e.Hinweis().String()
		},
	)
}

// this splits on right-expanded
func (s *SetPrefixNamed) Add(z id_set.Element) {
	es := kennung.Expanded(z.AkteEtiketten(), kennung.ExpanderRight)

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (s *SetPrefixNamed) addPair(e kennung.Etikett, z id_set.Element) {
	existing, ok := (*s)[e]

	if !ok {
		existing = makeMutableZettelLikeSet()
	}

	existing.Add(z)
	(*s)[e] = existing
}

// for all of the zettels, check for intersections with the passed in
// etikett, and if there is a prefix match, group it out the output set segments
// appropriately
func (a SetPrefixNamed) Subset(e kennung.Etikett) (out SetPrefixNamedSegments) {
	out.Ungrouped = makeMutableZettelLikeSet()
	out.Grouped = NewSetPrefixNamed()

	for e1, zSet := range a {
		zSet.Each(
			func(z id_set.Element) (err error) {
				intersection := z.AkteEtiketten().IntersectPrefixes(kennung.MakeEtikettSet(e))
				errors.Log().Printf("%s yields %s", e1, intersection)

				if intersection.Len() > 0 {
					for _, e2 := range intersection.Elements() {
						out.Grouped.addPair(e2, z)
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

func (s SetPrefixNamed) ToSetNamed() (out collections.MutableSet[id_set.Element]) {
	out = makeMutableZettelLikeSet()

	for _, zs := range s {
		zs.Each(
			func(z id_set.Element) (err error) {
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
