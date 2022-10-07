package id_set

import (
	"fmt"
	"reflect"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/charlie/typ"
)

type Set struct {
	shas sha.Set
	ids  []id.Id
}

func Make(c int) Set {
	return Set{
		shas: sha.MakeSet(c),
		ids:  make([]id.Id, 0, c),
	}
}

func (s *Set) Add(ids ...id.Id) {
	for _, i := range ids {
		switch it := i.(type) {
		case *sha.Sha:
		case sha.Sha:
			s.shas.Add(it)

		default:
			s.ids = append(s.ids, it)
		}
	}
}

func (s *Set) Shas() sha.Set {
	return s.shas
}

func (s Set) String() string {
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) Len() int {
	return len(s.ids) + s.shas.Len()
}

func (s Set) Hinweisen() (hinweisen []hinweis.Hinweis) {
	hinweisen = make([]hinweis.Hinweis, 0, len(s.ids))

	val := reflect.ValueOf(&hinweis.Hinweis{})
	t := val.Type()

	targetType := t.Elem()

	for _, i1 := range s.ids {
		if reflect.TypeOf(i1).AssignableTo(targetType) {
			hinweisen = append(hinweisen, i1.(hinweis.Hinweis))
		}
	}

	return
}

func (s Set) Timestamps() (timestamps []ts.Time) {
	timestamps = make([]ts.Time, 0, len(s.ids))

	val := reflect.ValueOf(&ts.Time{})
	t := val.Type()

	targetType := t.Elem()

	for _, i1 := range s.ids {
		if reflect.TypeOf(i1).AssignableTo(targetType) {
			timestamps = append(timestamps, i1.(ts.Time))
		}
	}

	return
}

func (s Set) Typen() (typen []typ.Typ) {
	typen = make([]typ.Typ, 0, len(s.ids))

	val := reflect.ValueOf(&typ.Typ{})
	t := val.Type()

	targetType := t.Elem()

	for _, i1 := range s.ids {
		if reflect.TypeOf(i1).AssignableTo(targetType) {
			typen = append(typen, i1.(typ.Typ))
		}
	}

	return
}

func (s Set) AnyShasOrHinweisen() (ids []id.IdMitKorper) {
	hinweisen := s.Hinweisen()
	ids = make([]id.IdMitKorper, 0, s.shas.Len()+len(hinweisen))

	s.shas.Each(
		func(sh sha.Sha) (err error) {
			ids = append(ids, sh)
			return
		},
	)

	for _, h := range hinweisen {
		ids = append(ids, h)
	}

	return
}

func (s Set) Etiketten() (etiketten etikett.Set) {
	mes := etikett.MakeMutableSet()

	val := reflect.ValueOf(&etikett.Etikett{})
	t := val.Type()

	targetType := t.Elem()

	for _, i1 := range s.ids {
		if reflect.TypeOf(i1).AssignableTo(targetType) {
			mes.Add(i1.(etikett.Etikett))
		}
	}

	etiketten = mes.Copy()

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
