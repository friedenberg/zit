package pool

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type fakePool[T any, TPtr schnittstellen.Resetable[T]] struct{}

func MakeFakePool[T any, TPtr schnittstellen.Resetable[T]]() *fakePool[T, TPtr] {
	return &fakePool[T, TPtr]{}
}

func (ip fakePool[T, TPtr]) Get() TPtr {
	var t T
	return &t
}

func (ip fakePool[T, TPtr]) Put(i TPtr) (err error) {
	return
}
