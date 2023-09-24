package schnittstellen

type (
	AkteLike interface{}

	AktePtrLike interface {
		AkteLike
	}

	Akte[T any] interface {
		Objekte[T]
		AkteLike
	}

	AktePtr[T any] interface {
		ObjektePtr[T]
		Resetable[T]
		AktePtrLike
	}
)
