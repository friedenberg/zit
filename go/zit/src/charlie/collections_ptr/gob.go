package collections_ptr

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

func RegisterGobValue[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
) {
	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
	}

	gob.Register(keyer)

	RegisterGob[T, TPtr]()
}

func RegisterGob[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]]() {
	gob.Register(Set[T, TPtr]{})
	gob.Register(MutableSet[T, TPtr]{})
}
