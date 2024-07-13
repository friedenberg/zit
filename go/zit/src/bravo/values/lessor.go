package values

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

type ReverseLessor[T any, TPtr interfaces.Ptr[T]] struct {
	Inner interfaces.Lessor2[T, TPtr]
}

func (rl ReverseLessor[T, TPtr]) Less(a T, b T) bool {
	return rl.Inner.Less(b, a)
}

func (rl ReverseLessor[T, TPtr]) LessPtr(a TPtr, b TPtr) bool {
	return rl.Inner.LessPtr(b, a)
}
