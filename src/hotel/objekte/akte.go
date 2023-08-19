package objekte

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type (
	AkteLike interface{}

	AktePtrLike interface {
		AkteLike
	}

	Akte[T any] interface {
		schnittstellen.Objekte[T]
		AkteLike
	}

	AktePtr[T any] interface {
		schnittstellen.ObjektePtr[T]
		schnittstellen.Resetable[T]
		AktePtrLike
	}
)
