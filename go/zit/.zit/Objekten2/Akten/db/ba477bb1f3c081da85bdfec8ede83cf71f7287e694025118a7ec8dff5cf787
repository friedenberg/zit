package quiter

import "code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"

func AddClonePoolPtr[E any, EPtr interfaces.Ptr[E]](
	s interfaces.Adder[EPtr],
	p interfaces.Pool[E, EPtr],
	r interfaces.Resetter2[E, EPtr],
	b EPtr,
) (err error) {
	a := p.Get()
	r.ResetWith(a, b)
	return s.Add(a)
}

func MakeAddClonePoolPtrFunc[E any, EPtr interfaces.Ptr[E]](
	s interfaces.Adder[EPtr],
	p interfaces.Pool[E, EPtr],
	r interfaces.Resetter2[E, EPtr],
) interfaces.FuncIter[EPtr] {
	return MakeSyncSerializer(func(e EPtr) (err error) {
		return AddClonePoolPtr(s, p, r, e)
	})
}

func MakeAddClonerPtrFunc[E interfaces.Cloner[E]](
	s interfaces.Adder[E],
) interfaces.FuncIter[E] {
	return MakeSyncSerializer(func(e E) (err error) {
		return s.Add(e.Clone())
	})
}
