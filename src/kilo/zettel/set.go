package zettel

import (
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
)

type MutableSet struct {
	collections.MutableSet[*Transacted]
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
					sz.Sku.Sha,
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
		func(z *Transacted) (err error) {
			b.Add(*z)

			return
		},
	)

	return
}

func (s MutableSet) ToSliceHinweisen() (b []hinweis.Hinweis) {
	b = make([]hinweis.Hinweis, 0, s.Len())

	s.Each(
		func(z *Transacted) (err error) {
			b = append(b, z.Sku.Kennung)

			return
		},
	)

	return
}