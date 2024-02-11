package values

import "code.linenisgreat.com/zit-go/src/alfa/schnittstellen"

type ReverseLessor[T any, TPtr schnittstellen.Ptr[T]] struct {
	Inner schnittstellen.Lessor2[T, TPtr]
}

func (rl ReverseLessor[T, TPtr]) Less(a T, b T) bool {
	return rl.Inner.Less(b, a)
}

func (rl ReverseLessor[T, TPtr]) LessPtr(a TPtr, b TPtr) bool {
	return rl.Inner.LessPtr(b, a)
}
