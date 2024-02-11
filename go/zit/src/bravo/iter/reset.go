package iter

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

func ResetMap[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

func ResetMutableSetWithPool[E any, EPtr schnittstellen.Ptr[E]](
	s schnittstellen.MutableSetPtrLike[E, EPtr],
	p schnittstellen.Pool[E, EPtr],
) {
	s.EachPtr(p.Put)
	s.Reset()
}
