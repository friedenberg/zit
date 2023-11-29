package iter

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func PrintPointer[T any, TPtr schnittstellen.Ptr[T]](e TPtr) (err error) {
	// log.Debug().Caller(1, "%s -> %s", unsafe.Pointer(e), e)
	return
}
