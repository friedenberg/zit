package id_set

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/hinweis"
	"github.com/friedenberg/zit/src/golf/ts"
	"github.com/friedenberg/zit/src/sha_collections"
)

// TODO-P4 move to kennung
// TODO-P3 rewrite
type Set struct {
	Shas       sha_collections.MutableSet
	Etiketten  kennung.EtikettMutableSet
	Hinweisen  hinweis.MutableSet
	Typen      kennung.TypMutableSet
	Timestamps ts.MutableSet
	HasKonfig  bool
	ids        []gattung.IdLike
}

func Make(c int) Set {
	return Set{
		Timestamps: ts.MakeMutableSet(),
		Shas:       sha_collections.MakeMutableSet(),
		Etiketten:  kennung.MakeEtikettMutableSet(),
		Hinweisen:  hinweis.MakeMutableSet(),
		Typen:      kennung.MakeTypMutableSet(),
		ids:        make([]gattung.IdLike, 0, c),
	}
}

func (s *Set) Add(ids ...gattung.IdLike) {
	for _, i := range ids {
		switch it := i.(type) {
		case kennung.Etikett:
			s.Etiketten.Add(it)

		case sha.Sha:
			s.Shas.Add(it)

		case hinweis.Hinweis:
			s.Hinweisen.Add(it)

		case kennung.Typ:
			s.Typen.Add(it)

		case ts.Time:
			s.Timestamps.Add(it)

		case kennung.Konfig:
			s.HasKonfig = true

		default:
			s.ids = append(s.ids, it)
		}
	}
}

func (s Set) String() string {
	errors.TodoP4("improve the string creation method")
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) Len() int {
	k := 0

	if s.HasKonfig {
		k = 1
	}

	return s.Shas.Len() + s.Etiketten.Len() + s.Hinweisen.Len() + s.Typen.Len() + s.Timestamps.Len() + k
}

func (s Set) AnyShasOrHinweisen() (ids []gattung.IdMitKorper) {
	ids = make([]gattung.IdMitKorper, 0, s.Shas.Len()+s.Hinweisen.Len())

	s.Shas.Each(
		func(sh sha.Sha) (err error) {
			ids = append(ids, sh)

			return
		},
	)

	s.Hinweisen.Each(
		func(h hinweis.Hinweis) (err error) {
			ids = append(ids, h)

			return
		},
	)

	return
}

func (s Set) AnyShaOrHinweis() (i1 gattung.IdMitKorper, ok bool) {
	ids := s.AnyShasOrHinweisen()

	if len(ids) > 0 {
		i1 = ids[0]
		ok = true
	}

	return
}
