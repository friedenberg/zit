package id_set

import (
	"fmt"

	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/typ"
)

type Set struct {
	shas       sha.MutableSet
	etiketten  etikett.MutableSet
	hinweisen  hinweis.MutableSet
	typen      typ.MutableSet
	timestamps ts.MutableSet
	hasKonfig  bool
	ids        []id.Id
}

func Make(c int) Set {
	return Set{
		shas:       sha.MakeMutableSet(),
		etiketten:  etikett.MakeMutableSet(),
		hinweisen:  hinweis.MakeMutableSet(),
		typen:      typ.MakeMutableSet(),
		timestamps: ts.MakeMutableSet(),
		ids:        make([]id.Id, 0, c),
	}
}

func (s *Set) Add(ids ...id.Id) {
	for _, i := range ids {
		switch it := i.(type) {
		case etikett.Etikett:
			s.etiketten.Add(it)

		case sha.Sha:
			s.shas.Add(it)

		case hinweis.Hinweis:
			s.hinweisen.Add(it)

		case typ.Typ:
			s.typen.Add(it)

		case ts.Time:
			s.timestamps.Add(it)

		case konfig.Id:
			s.hasKonfig = true

		default:
			s.ids = append(s.ids, it)
		}
	}
}

func (s *Set) Shas() (shas sha.Set) {
	return s.shas.Copy()
}

func (s Set) String() string {
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) Len() int {
	return len(s.ids) + s.shas.Len()
}

func (s Set) Hinweisen() (hinweisen hinweis.Set) {
	hinweisen = s.hinweisen.Copy()

	return
}

func (s Set) Timestamps() (timestamps []ts.Time) {
	timestamps = s.timestamps.Elements()

	return
}

func (s Set) Typen() (typen []typ.Typ) {
	typen = s.typen.Elements()

	return
}

func (s Set) HasKonfig() (ok bool) {
	ok = s.hasKonfig

	return
}

func (s Set) AnyShasOrHinweisen() (ids []id.IdMitKorper) {
	hinweisen := s.Hinweisen()
	ids = make([]id.IdMitKorper, 0, s.shas.Len()+hinweisen.Len())

	s.shas.Each(
		func(sh sha.Sha) (err error) {
			ids = append(ids, sh)

			return
		},
	)

	hinweisen.Each(
		func(h hinweis.Hinweis) (err error) {
			ids = append(ids, h)

			return
		},
	)

	return
}

func (s Set) Etiketten() (etiketten etikett.Set) {
	etiketten = s.etiketten.Copy()

	return
}

func (s Set) AnyShaOrHinweis() (i1 id.IdMitKorper, ok bool) {
	ids := s.AnyShasOrHinweisen()

	if len(ids) > 0 {
		i1 = ids[0]
		ok = true
	}

	return
}