package schnittstellen

type Ptr[T any] interface {
	*T
}

// TODO-P4 rename to Equaller
type Equatable[T any] interface {
	Equals(T) bool
}

type Flusher interface {
	Flush() error
}

type Resetter interface {
	Reset() error
}

// TODO-P4 rename to Resetter
type Resetable[T any] interface {
	Ptr[T]
	ResetWith(T)
	Reset()
}

type ResetWither[T any] interface {
	Ptr[T]
	ResetWith(T)
}

type Keyer[T any] interface {
	Key(T) string
}

type KeyPtrer[T any, T1 Ptr[T]] interface {
	Key(T1) string
}
