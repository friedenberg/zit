package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type SetPrefixNamed map[kennung.Etikett]schnittstellen.MutableSet[sku.WithKennungInterface]

type SetPrefixNamedSegments struct {
	Ungrouped schnittstellen.MutableSet[sku.WithKennungInterface]
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	errors.TodoP4("rename set prefix named")
	s := make(SetPrefixNamed)
	return &s
}

func makeMutableZettelLikeSet() schnittstellen.MutableSet[sku.WithKennungInterface] {
	return collections.MakeMutableSet(
		func(e sku.WithKennungInterface) string {
			return e.GetKennungLike().String()
		},
	)
}

// this splits on right-expanded
func (s *SetPrefixNamed) Add(z sku.WithKennungInterface) {
	es := kennung.Expanded(z.Metadatei.GetEtiketten(), kennung.ExpanderRight)

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (s *SetPrefixNamed) addPair(
	e kennung.Etikett,
	z sku.WithKennungInterface,
) {
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
			func(z sku.WithKennungInterface) (err error) {
				intersection := kennung.IntersectPrefixes(
					z.Metadatei.GetEtiketten(),
					kennung.MakeEtikettSet(e),
				)

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

func (s SetPrefixNamed) ToSetNamed() (out schnittstellen.MutableSet[sku.WithKennungInterface]) {
	out = makeMutableZettelLikeSet()

	for _, zs := range s {
		zs.Each(
			func(z sku.WithKennungInterface) (err error) {
				out.Add(z)

				return
			},
		)
	}

	return
}

func (a SetPrefixVerzeichnisse) Subset(
	e kennung.Etikett,
) (out SetPrefixVerzeichnisseSegments) {
	out.Ungrouped = MakeMutableSetMetadateiWithKennung()
	out.Grouped = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1.String() == "" {
			continue
		}

		zSet.Each(
			func(z sku.WithKennungInterface) (err error) {
				intersection := kennung.IntersectPrefixes(
					z.GetMetadatei().Etiketten,
					kennung.MakeEtikettSet(e),
				)
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
