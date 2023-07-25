package collections

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func RegisterGob[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]]() {
	gob.Register(&set[T]{})
	gob.Register(&mutableSet[T]{})
	gob.Register(MakeSetStringer[T, TPtr]())
	gob.Register(MakeMutableSetStringer[T]())
}
