package id_set

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
)

type Set struct {
	shas       sha.MutableSet
	etiketten  etikett.MutableSet
	hinweisen  hinweis.MutableSet
	typen      typ.MutableSet
	timestamps ts.MutableSet
	konfig     *konfig.Id
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
			s.konfig = &it

		default:
			s.ids = append(s.ids, it)
		}
	}
}

func (s *Set) Shas() (shas []sha.Sha) {
	return s.shas.Elements()
}

func (s Set) String() string {
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) Len() int {
	return len(s.ids) + s.shas.Len()
}

func (s Set) Hinweisen() (hinweisen []hinweis.Hinweis) {
	hinweisen = s.hinweisen.Elements()

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

func (s Set) Konfig() (ok bool) {
	ok = s.konfig != nil

	return
}

func (s Set) AnyShasOrHinweisen() (ids []id.IdMitKorper) {
	hinweisen := s.Hinweisen()
	ids = make([]id.IdMitKorper, 0, s.shas.Len()+len(hinweisen))

	s.shas.Each(
		func(sh sha.Sha) {
			ids = append(ids, sh)
		},
	)

	for _, h := range hinweisen {
		ids = append(ids, h)
	}

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
