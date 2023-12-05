package objekte_collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/expansion"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type SetPrefixVerzeichnisse struct {
	count    int
	innerMap map[string]MutableSetMetadateiWithKennung
}

type SetPrefixVerzeichnisseSegments struct {
	Ungrouped MutableSetMetadateiWithKennung
	Grouped   SetPrefixVerzeichnisse
}

func MakeSetPrefixVerzeichnisse(c int) (s SetPrefixVerzeichnisse) {
	s.innerMap = make(map[string]MutableSetMetadateiWithKennung, c)
	return s
}

func (s SetPrefixVerzeichnisse) Len() int {
	return s.count
}

// this splits on right-expanded
func (s *SetPrefixVerzeichnisse) Add(z *sku.Transacted) (err error) {
	es := kennung.Expanded(
		z.GetMetadatei().Verzeichnisse.GetImplicitEtiketten(),
		expansion.ExpanderRight,
	).CloneMutableSetPtrLike()

	if err = z.GetMetadatei().Verzeichnisse.GetExpandedEtiketten().EachPtr(
		es.AddPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if es.Len() == 0 {
		es = kennung.MakeEtikettMutableSet(kennung.Etikett{})
	}

	if err = es.Each(
		func(e kennung.Etikett) (err error) {
			s.addPair(e.String(), z)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a SetPrefixVerzeichnisse) Subtract(
	b MutableSetMetadateiWithKennung,
) (c SetPrefixVerzeichnisse) {
	c = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z *sku.Transacted) (err error) {
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
	e string,
	z *sku.Transacted,
) {
	s.count += 1

	existing, ok := s.innerMap[e]

	if !ok {
		existing = MakeMutableSetMetadateiWithKennung()
	}

	existing.Add(z)
	s.innerMap[e] = existing
}

func (a SetPrefixVerzeichnisse) Each(
	f func(kennung.Etikett, MutableSetMetadateiWithKennung) error,
) (err error) {
	for e, ssz := range a.innerMap {
		var e1 kennung.Etikett

		if e != "" {
			e1 = kennung.MustEtikett(e)
		}

		if err = f(e1, ssz); err != nil {
			if iter.IsStopIteration(err) {
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
	f func(kennung.Etikett, *sku.Transacted) error,
) error {
	return a.Each(
		func(e kennung.Etikett, st MutableSetMetadateiWithKennung) (err error) {
			err = st.Each(
				func(z *sku.Transacted) (err error) {
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
