package collections

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func RegisterGob[T schnittstellen.ValueLike]() {
	gob.Register(MakeSetStringer[T]())
	gob.Register(MakeMutableSetStringer[T]())
}
