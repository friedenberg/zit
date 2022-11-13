package id_set

import (
	"reflect"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/id"
)

type ProtoId struct {
	id.MutableId
	Expand func(string) (string, error)
}

type protoId struct {
	ProtoId
	reflect.Type
}

func makeProtoId(i ProtoId) protoId {
	return protoId{
		ProtoId: i,
		Type:    reflect.TypeOf(i.MutableId), // this type of this variable is reflect.Type
	}
}

func (pid protoId) String() string {
	return pid.Type.Name()
}

func (pid protoId) Make(v string) (i id.Id, err error) {
	if pid.Expand != nil {
		if v, err = pid.Expand(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	idPointer := reflect.New(pid.Type.Elem()) // this type of this variable is reflect.Value.
	idInterface := idPointer.Interface()      // this type of this variable is interface{}
	id2 := idInterface.(id.MutableId)

	if err = id2.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	i = reflect.ValueOf(id2).Elem().Interface().(id.Id)

	return
}
