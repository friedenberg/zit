package id_set

import (
	"fmt"
	"reflect"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/id"
)

var anyShaOrHinweisTypes []id.Id

func init() {
	anyShaOrHinweisTypes = []id.Id{
		&sha.Sha{},
		&hinweis.Hinweis{},
		&hinweis.HinweisWithIndex{},
	}
}

type Set struct {
	ids []id.Id
}

func (s Set) String() string {
	return fmt.Sprintf("%#v", s.ids)
}

func (s Set) Any(is ...id.Id) (i1 id.Id, ok bool) {
	for _, i := range is {
		val := reflect.ValueOf(i)
		typ := val.Type()

		targetType := typ.Elem()

		for _, i1 = range s.ids {
			if reflect.TypeOf(i1).AssignableTo(targetType) {
				ok = true
				return
			}
		}
	}

	return
}

func (s Set) AnyShaOrHinweis() (i1 id.Id, ok bool) {
	return s.Any(anyShaOrHinweisTypes...)
}
