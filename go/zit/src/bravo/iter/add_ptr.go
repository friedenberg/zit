package iter

import "code.linenisgreat.com/zit/src/alfa/schnittstellen"

func AddClonePoolPtr[E any, EPtr schnittstellen.Ptr[E]](
	s schnittstellen.Adder[EPtr],
	p schnittstellen.Pool[E, EPtr],
	r schnittstellen.Resetter2[E, EPtr],
	b EPtr,
) (err error) {
	a := p.Get()
	r.ResetWith(a, b)
	return s.Add(a)
}

func MakeAddClonePoolPtrFunc[E any, EPtr schnittstellen.Ptr[E]](
	s schnittstellen.Adder[EPtr],
	p schnittstellen.Pool[E, EPtr],
	r schnittstellen.Resetter2[E, EPtr],
) schnittstellen.FuncIter[EPtr] {
	return MakeSyncSerializer(func(e EPtr) (err error) {
		return AddClonePoolPtr(s, p, r, e)
	})
}
