package collections_value

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func RegisterGob[T schnittstellen.Element, TPtr schnittstellen.ElementPtr[T]]() {
	gob.Register(&Set[T]{})
	gob.Register(&MutableSet[T]{})
}
