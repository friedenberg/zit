package objekte_collections

import (
	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/expansion"
	"code.linenisgreat.com/zit-go/src/bravo/iter"
	"code.linenisgreat.com/zit-go/src/charlie/collections_value"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type SetPrefixNamed map[string]schnittstellen.MutableSetLike[*sku.Transacted]

type SetPrefixNamedSegments struct {
	Ungrouped schnittstellen.MutableSetLike[*sku.Transacted]
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	errors.TodoP4("rename set prefix named")
	s := make(SetPrefixNamed)
	return &s
}

func makeMutableZettelLikeSet() schnittstellen.MutableSetLike[*sku.Transacted] {
	return collections_value.MakeMutableSet[*sku.Transacted](
		KennungKeyer{},
	)
}

// this splits on right-expanded
func (s *SetPrefixNamed) Add(z *sku.Transacted) {
	es := kennung.Expanded(
		z.Metadatei.Verzeichnisse.GetImplicitEtiketten(),
		expansion.ExpanderRight,
	).CloneMutableSetPtrLike()

	var err error

	err = z.Metadatei.Verzeichnisse.GetExpandedEtiketten().EachPtr(es.AddPtr)
	errors.PanicIfError(err)

	err = es.Each(
		func(e kennung.Etikett) (err error) {
			s.addPair(e, z)
			return
		},
	)

	errors.PanicIfError(err)
}

func (s *SetPrefixNamed) addPair(
	e kennung.Etikett,
	z *sku.Transacted,
) {
	existing, ok := (*s)[e.String()]

	if !ok {
		existing = makeMutableZettelLikeSet()
	}

	existing.Add(z)
	(*s)[e.String()] = existing
}

func allEtiketten(z *sku.Transacted) (es kennung.EtikettMutableSet, err error) {
	es = kennung.MakeEtikettMutableSet()

	if err = z.Metadatei.Verzeichnisse.GetImplicitEtiketten().EachPtr(es.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = z.Metadatei.GetEtiketten().EachPtr(es.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// for all of the zettels, check for intersections with the passed in
// etikett, and if there is a prefix match, group it out the output set segments
// appropriately
func (a SetPrefixNamed) Subset(e kennung.Etikett) (out SetPrefixNamedSegments) {
	out.Ungrouped = makeMutableZettelLikeSet()
	out.Grouped = NewSetPrefixNamed()

	for _, zSet := range a {
		zSet.Each(
			func(z *sku.Transacted) (err error) {
				var es kennung.EtikettSet

				if es, err = allEtiketten(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				intersection := kennung.IntersectPrefixes(
					es,
					kennung.MakeEtikettSet(e),
				)

				exactMatch := intersection.Len() == 1 &&
					intersection.Any().Equals(e)

				if intersection.Len() > 0 && !exactMatch {
					for _, e2 := range iter.Elements[kennung.Etikett](intersection) {
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

func (s SetPrefixNamed) ToSetNamed() (out schnittstellen.MutableSetLike[*sku.Transacted]) {
	out = makeMutableZettelLikeSet()

	for _, zs := range s {
		zs.Each(
			func(z *sku.Transacted) (err error) {
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
		if e1 == "" {
			continue
		}

		zSet.Each(
			func(z *sku.Transacted) (err error) {
				var es kennung.EtikettSet

				if es, err = allEtiketten(z); err != nil {
					err = errors.Wrap(err)
					return
				}

				intersection := kennung.IntersectPrefixes(
					es,
					kennung.MakeEtikettSet(e),
				)

				exactMatch := intersection.Len() == 1 &&
					intersection.Any().Equals(e)

				if intersection.Len() > 0 && !exactMatch {
					for _, e2 := range iter.Elements[kennung.Etikett](intersection) {
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
