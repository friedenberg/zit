package collections_ptr

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func RegisterGob[T schnittstellen.Element, TPtr schnittstellen.ElementPtr[T]]() {
	gob.Register(&Set[T, TPtr]{})
	gob.Register(&MutableSet[T, TPtr]{})
}
