package kennung

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/sha_collections"
	"github.com/friedenberg/zit/src/echo/ts"
)

// TODO-P4 move to kennung
// TODO-P3 rewrite
type Set struct {
	Shas       sha_collections.MutableSet
	Etiketten  EtikettMutableSet
	Hinweisen  HinweisMutableSet
	Typen      TypMutableSet
	Timestamps ts.MutableSet
	HasKonfig  bool
	Sigil      Sigil
	ids        []schnittstellen.Value
}

func Make(c int) Set {
	return Set{
		Timestamps: ts.MakeMutableSet(),
		Shas:       sha_collections.MakeMutableSet(),
		Etiketten:  MakeEtikettMutableSet(),
		Hinweisen:  MakeHinweisMutableSet(),
		Typen:      MakeTypMutableSet(),
		ids:        make([]schnittstellen.Value, 0, c),
	}
}

func (s *Set) Add(ids ...schnittstellen.Value) {
	for _, i := range ids {
		switch it := i.(type) {
		case Etikett:
			s.Etiketten.Add(it)

		case sha.Sha:
			s.Shas.Add(it)

		case Hinweis:
			s.Hinweisen.Add(it)

		case Typ:
			s.Typen.Add(it)

		case ts.Time:
			s.Timestamps.Add(it)

		case Konfig:
			s.HasKonfig = true

		case Sigil:
			s.Sigil.Add(it)

		default:
			s.ids = append(s.ids, it)
		}
	}
}

func (s Set) String() string {
	errors.TodoP4("improve the string creation method")
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) OnlySingleHinweis() (h Hinweis, ok bool) {
	h = s.Hinweisen.Any()
	ok = s.Len() == 1 && s.Hinweisen.Len() == 1 && !h.GetSigil().IncludesHistory()

	return
}

func (s Set) Len() int {
	k := 0

	if s.HasKonfig {
		k = 1
	}

	return s.Shas.Len() + s.Etiketten.Len() + s.Hinweisen.Len() + s.Typen.Len() + s.Timestamps.Len() + k
}

func (s Set) AnyShasOrHinweisen() (ids []schnittstellen.Korper) {
	ids = make([]schnittstellen.Korper, 0, s.Shas.Len()+s.Hinweisen.Len())

	s.Shas.Each(
		func(sh sha.Sha) (err error) {
			ids = append(ids, sh)

			return
		},
	)

	s.Hinweisen.Each(
		func(h Hinweis) (err error) {
			ids = append(ids, h)

			return
		},
	)

	return
}

func (s Set) AnyShaOrHinweis() (i1 schnittstellen.Korper, ok bool) {
	ids := s.AnyShasOrHinweisen()

	if len(ids) > 0 {
		i1 = ids[0]
		ok = true
	}

	return
}