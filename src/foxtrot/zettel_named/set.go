package zettel_named

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/hinweis"
)

// TODO deprecate and move to MutableSet
type SetNamed struct {
	collections.MutableSetLike[Zettel]
}

func NewSetNamed() *SetNamed {
	return &SetNamed{
		MutableSetLike: collections.MakeMutableSetGeneric(
			func(sz Zettel) string {
				// if sz == nil {
				// 	return ""
				// }

				return collections.MakeKey(
					sz.Hinweis,
				)
			},
		),
	}
}

func (s SetNamed) Hinweisen() (h []hinweis.Hinweis) {
	h = make([]hinweis.Hinweis, 0, s.Len())

	s.Each(
		func(z Zettel) (err error) {
			h = append(h, z.Hinweis)

			return
		},
	)

	return
}

func (s SetNamed) HinweisStrings() (h []string) {
	h = make([]string, 0, s.Len())

	s.Each(
		func(z Zettel) (err error) {
			h = append(h, z.Hinweis.String())

			return
		},
	)

	return
}
