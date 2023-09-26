package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type SetPrefixNamed map[kennung.Etikett]schnittstellen.MutableSetLike[sku.SkuLikePtr]

type SetPrefixNamedSegments struct {
	Ungrouped schnittstellen.MutableSetLike[sku.SkuLikePtr]
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	errors.TodoP4("rename set prefix named")
	s := make(SetPrefixNamed)
	return &s
}

func makeMutableZettelLikeSet() schnittstellen.MutableSetLike[sku.SkuLikePtr] {
	return collections_value.MakeMutableSet[sku.SkuLikePtr](
		KennungKeyer{},
	)
}

// this splits on right-expanded
func (s *SetPrefixNamed) Add(z sku.SkuLikePtr) {
	es := kennung.Expanded(
		z.GetMetadatei().GetEtiketten(),
		kennung.ExpanderRight,
	)

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (s *SetPrefixNamed) addPair(
	e kennung.Etikett,
	z sku.SkuLikePtr,
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

	for _, zSet := range a {
		zSet.Each(
			func(z sku.SkuLikePtr) (err error) {
				intersection := kennung.IntersectPrefixes(
					z.GetMetadatei().GetEtiketten(),
					kennung.MakeEtikettSet(e),
				)

				exactMatch := intersection.Len() == 1 &&
					intersection.Any().Equals(e)

				if intersection.Len() > 0 && !exactMatch {
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

func (s SetPrefixNamed) ToSetNamed() (out schnittstellen.MutableSetLike[sku.SkuLikePtr]) {
	out = makeMutableZettelLikeSet()

	for _, zs := range s {
		zs.Each(
			func(z sku.SkuLikePtr) (err error) {
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
			func(z sku.SkuLikePtr) (err error) {
				intersection := kennung.IntersectPrefixes(
					z.GetMetadatei().Etiketten,
					kennung.MakeEtikettSet(e),
				)

				exactMatch := intersection.Len() == 1 &&
					intersection.Any().Equals(e)

				if intersection.Len() > 0 && !exactMatch {
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
