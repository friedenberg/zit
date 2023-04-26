package objekte

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type Objekte[T any] interface {
	schnittstellen.Objekte[T]
	ObjekteLike
}

type ObjektePtr[T any] interface {
	schnittstellen.ObjektePtr[T]
	schnittstellen.Resetable[T]
	ObjektePtrLike
}
