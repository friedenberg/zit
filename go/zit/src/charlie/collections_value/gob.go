package collections_value

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

func RegisterGobValue[T schnittstellen.ValueLike](
	keyer schnittstellen.StringKeyer[T],
) {
	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
	}

	gob.Register(keyer)

	RegisterGob[T]()
}

func RegisterGob[T schnittstellen.ValueLike]() {
	gob.Register(Set[T]{})
	gob.Register(MutableSet[T]{})
}
