package iter

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type StringerKeyer[
	T interfaces.Stringer,
] struct{}

func (sk StringerKeyer[T]) RegisterGob() StringerKeyer[T] {
	gob.Register(sk)
	return sk
}

func (sk StringerKeyer[T]) GetKey(e T) string {
	return e.String()
}

type StringerKeyerPtr[
	T interfaces.Stringer,
	TPtr interface {
		interfaces.Ptr[T]
		interfaces.Stringer
	},
] struct{}

func (sk StringerKeyerPtr[T, TPtr]) RegisterGob() StringerKeyerPtr[T, TPtr] {
	gob.Register(sk)
	return sk
}

func (sk StringerKeyerPtr[T, TPtr]) GetKey(e T) string {
	return e.String()
}

func (sk StringerKeyerPtr[T, TPtr]) GetKeyPtr(e TPtr) string {
	if e == nil {
		return ""
	}

	return e.String()
}
