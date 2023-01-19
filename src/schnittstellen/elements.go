package schnittstellen

type Value interface {
	Stringer
}

type ValuePtr[T Value] interface {
	Ptr[T]
	Setter
}

type IdLike interface {
	GattungGetter
	Value
}

type Id[T Value] interface {
	Value
	GattungGetter
	Equatable[T]
}

type IdPtr[T Value] interface {
	Id[T]
	ValuePtr[T]
	Resetable[T]
}
