package id_set

import (
	"reflect"

	"github.com/friedenberg/zit/src/delta/id"
)

type ProtoSet struct {
	types map[string]reflect.Type
}

func MakeProtoSet(types ...id.MutableId) (ps ProtoSet) {
	ps.types = make(map[string]reflect.Type, len(types))

	for _, t := range types {
		idType := reflect.TypeOf(t) // this type of this variable is reflect.Type
		ps.types[idType.Name()] = idType
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
		idPointer := reflect.New(t.Elem())   // this type of this variable is reflect.Value.
		idInterface := idPointer.Interface() // this type of this variable is interface{}
		id2 := idInterface.(id.MutableId)

		if err := id2.Set(v); err == nil {
			id := reflect.ValueOf(id2).Elem().Interface().(id.Id)
			s.ids = append(s.ids, id)
		}
	}

	return
}
