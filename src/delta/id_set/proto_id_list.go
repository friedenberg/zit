package id_set

import (
	"reflect"

	"github.com/friedenberg/zit/src/alfa/errors"
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

func (ps ProtoIdList) Make(vs ...string) (s Set, err error) {
	s = Make(len(vs))

	for _, v := range vs {
		for _, t := range ps.types {
			var i id.Id

			if i, err = t.Make(v); err == nil {
				s.Add(i)
				break
			}
		}

		if err != nil {
			err = errors.Errorf("no proto id was able to parse: %s", v)
			return
		}
	}

	return
}
