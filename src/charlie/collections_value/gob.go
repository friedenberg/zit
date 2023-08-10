package collections_value

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

func RegisterGobValue[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyer[T],
) {
	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
	}

	gob.Register(keyer)

	gob.Register(&Set[T]{})
	gob.Register(MutableSet[T]{})
}
