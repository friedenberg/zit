package interfaces

type Element interface {
	EqualsAny(any) bool
}

type ElementPtr[T any] interface {
	Ptr[T]
	Element
}

type ValueLike interface {
	Stringer
	Element
}

type ValuePtrLike interface {
	ValueLike
	Setter
}

type Value[T any] interface {
	ValueLike
	Equatable[T]
}

type ValuePtr[T any] interface {
	ValueLike
	// Value[T]
	Ptr[T]
}

type Id[T any] interface {
	Value[T]
	GenreGetter
}

type IdPtr[T any] interface {
	Id[T]
	ValuePtr[T]
	Resetable[T]
}

// TODO-P2 remove
type Lessor[T any] interface {
	Less(T) bool
}

type Lessor3[T any] interface {
	Less(T, T) bool
}

// TODO-P2 rename
type Lessor2[T any, TPtr Ptr[T]] interface {
	Lessor3[T]
	LessPtr(TPtr, TPtr) bool
}

// TODO-P2 rename
type Equaler1[T any] interface {
	Equals(T, T) bool
}

// TODO-P2 rename
type Equaler[T any, TPtr Ptr[T]] interface {
	Equaler1[T]
	EqualsPtr(TPtr, TPtr) bool
}

type Resetter2[T any, TPtr Ptr[T]] interface {
	Reset(TPtr)
	ResetWith(TPtr, TPtr)
}

type Resetter3[T any] interface {
	Reset(T)
	ResetWith(T, T)
}
