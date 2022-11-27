package zettel_transacted

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/hinweis"
)

type MutableSet struct {
	collections.MutableSet[*Zettel]
}

func MakeMutableSetUnique(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Zettel) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Kopf,
					sz.Mutter[0],
					sz.Mutter[1],
					sz.Schwanz,
					sz.TransaktionIndex,
					sz.Named.Kennung,
					sz.Named.Stored.Sha,
				)
			},
		),
	}
}

func MakeMutableSetHinweis(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Zettel) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Named.Kennung,
				)
			},
		),
	}
}

// func (m Set) Filter(w Writer) (err error) {
// 	for k, sz := range m.innerMap {
// 		if err = w.WriteZettelTransacted(&sz); err != nil {
// 			if errors.IsEOF(err) {
// 				err = nil
// 				delete(m.innerMap, k)
// 			} else {
// 				err = errors.Wrap(err)
// 				return
// 			}
// 		}
// 	}

// 	return
// }

func (s MutableSet) ToSetPrefixTransacted() (b SetPrefixTransacted) {
	b = MakeSetPrefixTransacted(s.Len())

	s.Each(
		func(z *Zettel) (err error) {
			b.Add(*z)

			return
		},
	)

	return
}

func (s MutableSet) ToSliceHinweisen() (b []hinweis.Hinweis) {
	b = make([]hinweis.Hinweis, 0, s.Len())

	s.Each(
		func(z *Zettel) (err error) {
			b = append(b, z.Named.Kennung)

			return
		},
	)

	return
}
