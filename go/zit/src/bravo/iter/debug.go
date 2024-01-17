package iter

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func PrintPointer[T any, TPtr schnittstellen.Ptr[T]](e TPtr) (err error) {
	// log.Debug().Caller(1, "%s -> %s", unsafe.Pointer(e), e)
	return
}

func MakeIterDebug[T any](f schnittstellen.FuncIter[T]) schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		// log.Debug().FunctionName(2)
		return f(e)
	}
}
