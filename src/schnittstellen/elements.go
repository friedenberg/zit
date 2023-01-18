package schnittstellen

type Value interface {
	Stringer
}

type ValuePtr[T Value] interface {
	Ptr[T]
	Setter
}
