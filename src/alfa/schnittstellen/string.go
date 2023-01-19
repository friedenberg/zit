package schnittstellen

type FuncSetString func(string) error
type FuncAbbreviateValue func(Value) (string, error)
type FuncAbbreviateKorper func(Korper) (string, error)

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
