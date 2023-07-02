package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type SetPrefixVerzeichnisse struct {
	count    int
	innerMap map[kennung.Etikett]MutableSetMetadateiWithKennung
}

type SetPrefixVerzeichnisseSegments struct {
	Ungrouped MutableSetMetadateiWithKennung
	Grouped   SetPrefixVerzeichnisse
}

func MakeSetPrefixVerzeichnisse(c int) (s SetPrefixVerzeichnisse) {
	s.innerMap = make(map[kennung.Etikett]MutableSetMetadateiWithKennung, c)
	return s
}

func (s SetPrefixVerzeichnisse) Len() int {
	return s.count
}

// this splits on right-expanded
func (s *SetPrefixVerzeichnisse) Add(z sku.WithKennungInterface) (err error) {
	es := kennung.Expanded(z.Metadatei.GetEtiketten(), kennung.ExpanderRight)

	if es.Len() == 0 {
		es = kennung.MakeEtikettSet(kennung.Etikett{})
	}

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}

	return
}

func (a SetPrefixVerzeichnisse) Subtract(b MutableSetMetadateiWithKennung) (c SetPrefixVerzeichnisse) {
	c = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z sku.WithKennungInterface) (err error) {
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

func (s *SetPrefixVerzeichnisse) addPair(
	e kennung.Etikett,
	z sku.WithKennungInterface,
) {
	s.count += 1

	existing, ok := s.innerMap[e]

	if !ok {
		existing = MakeMutableSetMetadateiWithKennung()
	}

	existing.Add(z)
	s.innerMap[e] = existing
}

func (a SetPrefixVerzeichnisse) Each(f func(kennung.Etikett, MutableSetMetadateiWithKennung) error) (err error) {
	for e, ssz := range a.innerMap {
		if err = f(e, ssz); err != nil {
			if collections.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a SetPrefixVerzeichnisse) EachZettel(
	f func(kennung.Etikett, sku.WithKennungInterface) error,
) error {
	return a.Each(
		func(e kennung.Etikett, st MutableSetMetadateiWithKennung) (err error) {
			st.Each(
				func(z sku.WithKennungInterface) (err error) {
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

func (s SetPrefixVerzeichnisse) ToSet() (out MutableSetMetadateiWithKennung) {
	out = MakeMutableSetMetadateiWithKennung()

	for _, zs := range s.innerMap {
		zs.Each(out.Add)
	}

	return
}
