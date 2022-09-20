package id_set

import (
	"reflect"

	"github.com/friedenberg/zit/src/charlie/id"
)

type ProtoIdList struct {
	types []protoId
}

func MakeProtoIdList(types ...ProtoId) (ps ProtoIdList) {
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

func (ps ProtoIdList) Len() int {
	return len(ps.types)
}

func (ps ProtoIdList) Make(vs ...string) (s Set) {
	s = Set{
		ids: make([]id.Id, 0, len(vs)),
	}

	for _, v := range vs {
		for _, t := range ps.types {
			var i id.Id
			var err error

			if i, err = t.Make(v); err == nil {
				s.ids = append(s.ids, i)
				break
			}
		}
	}

	return
}
