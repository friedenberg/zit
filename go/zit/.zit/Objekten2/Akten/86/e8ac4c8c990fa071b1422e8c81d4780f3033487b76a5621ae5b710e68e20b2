package collections_value

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
)

func RegisterGobValue[T interfaces.ValueLike](
	keyer interfaces.StringKeyer[T],
) {
	if keyer == nil {
		keyer = quiter.StringerKeyer[T]{}.RegisterGob()
	}

	gob.Register(keyer)

	RegisterGob[T]()
}

func RegisterGob[T interfaces.ValueLike]() {
	gob.Register(Set[T]{})
	gob.Register(MutableSet[T]{})
}
