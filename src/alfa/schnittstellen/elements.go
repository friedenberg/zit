package schnittstellen

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
	GattungGetter
}

type IdPtr[T any] interface {
	Id[T]
	ValuePtr[T]
	Resetable[T]
}

type Lessor[T any] interface {
	Less(T) bool
}

type Lessor2[T any, TPtr Ptr[T]] interface {
	Less(T, T) bool
	LessPtr(TPtr, TPtr) bool
}

type Equaler[T any, TPtr Ptr[T]] interface {
	Equals(T, T) bool
	EqualsPtr(TPtr, TPtr) bool
}
