package iter

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type StringerKeyer[
	T schnittstellen.Stringer,
	TPtr interface {
		schnittstellen.Ptr[T]
		schnittstellen.Stringer
	},
] struct{}

func (sk StringerKeyer[T, TPtr]) RegisterGob() StringerKeyer[T, TPtr] {
	gob.Register(sk)
	return sk
}

func (sk StringerKeyer[T, TPtr]) GetKey(e T) string {
	return e.String()
}

func (sk StringerKeyer[T, TPtr]) GetKeyPtr(e TPtr) string {
	if e == nil {
		return ""
	}

	return e.String()
}
