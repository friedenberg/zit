package id_set

import (
	"reflect"

	"github.com/friedenberg/zit/src/charlie/id"
)

type ProtoSet struct {
	types []protoId
}

func MakeProtoSet(types ...ProtoId) (ps ProtoSet) {
	ps.types = make([]protoId, len(types))

	for i, t := range types {
		pid := protoId{
			ProtoId: t,
			Type:    reflect.TypeOf(t.MutableId), // this type of this variable is reflect.Type
		}

		ps.types[i] = pid
	}

	return
}

func (ps ProtoSet) MakeMany(vs ...string) (ss []Set) {
	ss = make([]Set, len(vs))

	for i, v := range vs {
		ss[i] = ps.MakeOne(v)
	}

	return
}

func (ps ProtoSet) MakeOne(v string) (s Set) {
	for _, t := range ps.types {
		var i id.Id
		var err error

		if i, err = t.Make(v); err == nil {
			s.ids = append(s.ids, i)
			break
		}
	}

	return
}
