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

type Setter interface {
	Set(string) error
}

type SetterPtr[T any] interface {
	Ptr[T]
	Setter
}
