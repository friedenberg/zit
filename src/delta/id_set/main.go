package id_set

import (
	"fmt"
	"reflect"

	"github.com/friedenberg/zit/src/alfa/typ"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/charlie/id"
	"github.com/friedenberg/zit/src/charlie/ts"
)

type Set struct {
	ids []id.Id
}

func (s Set) String() string {
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) Len() int {
	return len(s.ids)
}

func (s Set) Shas() (shas []sha.Sha) {
	shas = make([]sha.Sha, 0, len(s.ids))

	val := reflect.ValueOf(&sha.Sha{})
	t := val.Type()

	targetType := t.Elem()

	for _, i1 := range s.ids {
		if reflect.TypeOf(i1).AssignableTo(targetType) {
			shas = append(shas, i1.(sha.Sha))
		}
	}

	return
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
	shas := s.Shas()
	hinweisen := s.Hinweisen()
	ids = make([]id.IdMitKorper, 0, len(shas)+len(hinweisen))

	for _, sh := range shas {
		ids = append(ids, sh)
	}

	for _, h := range hinweisen {
		ids = append(ids, h)
	}

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
