package schnittstellen

type (
	FuncSetString        func(string) error
	FuncString[T any]    func(T) string
	FuncAbbreviateValue  func(ValueLike) (string, error)
	FuncAbbreviateKorper func(Korper) (string, error)
)

type Stringer interface {
	String() string
}

type StringLenner interface {
	Stringer
	Lenner
}

type Setter interface {
	Set(string) error
}

type SetterPtr[T any] interface {
	Ptr[T]
	Setter
}

type StringSetterPtr[T any] interface {
	Stringer
	Ptr[T]
	Setter
}
