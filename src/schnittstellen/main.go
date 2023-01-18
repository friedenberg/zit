package schnittstellen

type Ptr[T any] interface {
	*T
}

type Equatable[T any] interface {
	Equals(T) bool
}

