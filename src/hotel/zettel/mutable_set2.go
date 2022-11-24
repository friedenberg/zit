package zettel

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/hinweis"
)

type NamedMutableSet struct {
	collections.MutableSetLike[*Named]
}

func MakeNamedMutableSet() NamedMutableSet {
	return NamedMutableSet{
		MutableSetLike: collections.MakeMutableSet(
			func(sz *Named) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Kennung,
				)
			},
		),
	}
}

func (s NamedMutableSet) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, s.Len())

	s.Each(
		func(z *Named) (err error) {
			h = append(h, z.Kennung)

			return
		},
	)

	return
}

func (s NamedMutableSet) HinweisStrings() (h []string) {
	h = make([]string, 0, s.Len())

	s.Each(
		func(z *Named) (err error) {
			h = append(h, z.Kennung.String())

			return
		},
	)

	return
}
