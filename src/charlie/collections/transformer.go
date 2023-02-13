package collections

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func MakeFuncTransformer[T any, T1 any](wf schnittstellen.FuncIter[T]) schnittstellen.FuncIter[T1] {
	return func(e T1) (err error) {
		if e1, ok := any(e).(T); ok {
			return wf(e1)
		}

		return
	}
}
