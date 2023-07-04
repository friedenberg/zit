package schnittstellen

type FuncError func() error

type Ptr[T any] interface {
	*T
}

type Equatable[T any] interface {
	Equals(T) bool
}

type Flusher interface {
	Flush() error
}

type Resetter interface {
	Reset()
}

type ResetterWithError interface {
	Reset() error
}

type Resetable[T any] interface {
	Ptr[T]
	ResetWith(T)
	Reset()
}

type Keyer[T any] interface {
	Key(T) string
}

type KeyPtrer[T any, T1 Ptr[T]] interface {
	Key(T1) string
}
