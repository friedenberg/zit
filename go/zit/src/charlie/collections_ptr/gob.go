package collections_ptr

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
)

func RegisterGobValue[T interfaces.ValueLike, TPtr interfaces.ValuePtr[T]](
	keyer interfaces.StringKeyerPtr[T, TPtr],
) {
	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
	}

	gob.Register(keyer)

	RegisterGob[T, TPtr]()
}

func RegisterGob[T interfaces.ValueLike, TPtr interfaces.ValuePtr[T]]() {
	gob.Register(Set[T, TPtr]{})
	gob.Register(MutableSet[T, TPtr]{})
}
