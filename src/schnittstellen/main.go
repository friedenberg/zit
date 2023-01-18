package schnittstellen

type Ptr[T any] interface {
	*T
}

type Equatable[T any] interface {
	Equals(T) bool
}

type Resetable[T any] interface {
	Ptr[T]
	ResetWith(T)
	Reset()
}

type ResetWither[T any] interface {
	Ptr[T]
	ResetWith(T)
}
