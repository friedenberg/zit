package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type HeapTransacted = collections.Heap[Transacted, *Transacted]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[Transacted, *Transacted]()
}

type MutableSet struct {
	schnittstellen.MutableSet[*Transacted]
}

func MakeMutableSetUnique(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Transacted) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Sku.Kopf,
					sz.Sku.Mutter[0],
					sz.Sku.Mutter[1],
					sz.Sku.Schwanz,
					sz.Sku.TransactionIndex,
					sz.Sku.Kennung,
					sz.Sku.ObjekteSha,
				)
			},
		),
	}
}

func MakeMutableSetHinweis(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Transacted) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Sku.Kennung,
				)
			},
		),
	}
}

func (s MutableSet) ToSetPrefixVerzeichnisse() (b SetPrefixVerzeichnisse) {
	b = MakeSetPrefixVerzeichnisse(s.Len())

	s.Each(
		func(z *Transacted) (err error) {
			b.Add(*z)

			return
		},
	)

	return
}

func (s MutableSet) ToSliceHinweisen() (b []kennung.Hinweis) {
	b = make([]kennung.Hinweis, 0, s.Len())

	s.Each(
		func(z *Transacted) (err error) {
			b = append(b, z.Sku.Kennung)

			return
		},
	)

	return
}

//   ____       _   ____            __ _
//  / ___|  ___| |_|  _ \ _ __ ___ / _(_)_  __
//  \___ \ / _ \ __| |_) | '__/ _ \ |_| \ \/ /
//   ___) |  __/ |_|  __/| | |  __/  _| |>  <
//  |____/ \___|\__|_|   |_|  \___|_| |_/_/\_\
//

type SetPrefixVerzeichnisse struct {
	count    int
	innerMap map[kennung.Etikett]MutableSet
}

type SetPrefixVerzeichnisseSegments struct {
	Ungrouped MutableSet
	Grouped   SetPrefixVerzeichnisse
}

func MakeSetPrefixVerzeichnisse(c int) (s SetPrefixVerzeichnisse) {
	s.innerMap = make(map[kennung.Etikett]MutableSet, c)
	return s
}

func (s SetPrefixVerzeichnisse) Len() int {
	return s.count
}

// this splits on right-expanded
func (s *SetPrefixVerzeichnisse) Add(z Transacted) {
	es := kennung.Expanded(z.Objekte.Metadatei.Etiketten, kennung.ExpanderRight)

	if es.Len() == 0 {
		es = kennung.MakeEtikettSet(kennung.Etikett{})
	}

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (a SetPrefixVerzeichnisse) Subtract(b MutableSet) (c SetPrefixVerzeichnisse) {
	c = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e, aSet := range a.innerMap {
		aSet.Each(
			func(z *Transacted) (err error) {
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

func (s *SetPrefixVerzeichnisse) addPair(e kennung.Etikett, z Transacted) {
	s.count += 1

	existing, ok := s.innerMap[e]

	if !ok {
		existing = MakeMutableSetUnique(1)
	}

	existing.Add(&z)
	s.innerMap[e] = existing
}

func (a SetPrefixVerzeichnisse) Each(f func(kennung.Etikett, MutableSet) error) (err error) {
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

func (a SetPrefixVerzeichnisse) EachZettel(f func(kennung.Etikett, Transacted) error) error {
	return a.Each(
		func(e kennung.Etikett, st MutableSet) (err error) {
			st.Each(
				func(z *Transacted) (err error) {
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
func (a SetPrefixVerzeichnisse) Subset(e kennung.Etikett) (out SetPrefixVerzeichnisseSegments) {
	out.Ungrouped = MakeMutableSetUnique(len(a.innerMap))
	out.Grouped = MakeSetPrefixVerzeichnisse(len(a.innerMap))

	for e1, zSet := range a.innerMap {
		if e1.String() == "" {
			continue
		}

		zSet.Each(
			func(z *Transacted) (err error) {
				intersection := kennung.IntersectPrefixes(
					z.Objekte.Metadatei.Etiketten,
					kennung.MakeEtikettSet(e),
				)
				errors.Log().Printf("%s yields %s", e1, intersection)

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

func (s SetPrefixVerzeichnisse) ToSet() (out MutableSet) {
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

//   ____       _   ____            __ _      _   _                          _
//  / ___|  ___| |_|  _ \ _ __ ___ / _(_)_  _| \ | | __ _ _ __ ___   ___  __| |
//  \___ \ / _ \ __| |_) | '__/ _ \ |_| \ \/ /  \| |/ _` | '_ ` _ \ / _ \/ _` |
//   ___) |  __/ |_|  __/| | |  __/  _| |>  <| |\  | (_| | | | | | |  __/ (_| |
//  |____/ \___|\__|_|   |_|  \___|_| |_/_/\_\_| \_|\__,_|_| |_| |_|\___|\__,_|
//

type SetPrefixNamed map[kennung.Etikett]schnittstellen.MutableSet[kennung.Matchable]

type SetPrefixNamedSegments struct {
	Ungrouped schnittstellen.MutableSet[kennung.Matchable]
	Grouped   *SetPrefixNamed
}

func NewSetPrefixNamed() *SetPrefixNamed {
	errors.TodoP4("rename set prefix named")
	s := make(SetPrefixNamed)
	return &s
}

func makeMutableZettelLikeSet() schnittstellen.MutableSet[kennung.Matchable] {
	return collections.MakeMutableSet(
		func(e kennung.Matchable) string {
			if e == nil {
				return ""
			}

			return e.GetIdLike().String()
		},
	)
}

// this splits on right-expanded
func (s *SetPrefixNamed) Add(z kennung.Matchable) {
	es := kennung.Expanded(z.GetEtiketten(), kennung.ExpanderRight)

	for _, e := range es.Elements() {
		s.addPair(e, z)
	}
}

func (s *SetPrefixNamed) addPair(e kennung.Etikett, z kennung.Matchable) {
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
			func(z kennung.Matchable) (err error) {
				intersection := kennung.IntersectPrefixes(
					z.GetEtiketten(),
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

func (s SetPrefixNamed) ToSetNamed() (out schnittstellen.MutableSet[kennung.Matchable]) {
	out = makeMutableZettelLikeSet()

	for _, zs := range s {
		zs.Each(
			func(z kennung.Matchable) (err error) {
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
