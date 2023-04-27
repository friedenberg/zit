package objekte

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type Akte[T any] interface {
	schnittstellen.Objekte[T]
	ObjekteLike
}

type AktePtr[T any] interface {
	schnittstellen.ObjektePtr[T]
	schnittstellen.Resetable[T]
	ObjektePtrLike
}
