package schnittstellen

type Element interface {
	EqualsAny(any) bool
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
	Value[T]
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
