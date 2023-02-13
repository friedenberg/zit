package schnittstellen

import "flag"

type (
	FuncIter[T any]        func(T) error
	FuncIterKey            func(string) error
	FuncIterWithKey[T any] func(string, T) error
)

type SetLike[T any] interface {
	Len() int
	Key(T) string
	Get(string) (T, bool)
	ContainsKey(string) bool
	Contains(T) bool
	Each(FuncIter[T]) error
	EachPtr(FuncIter[*T]) error
	EachKey(FuncIterKey) error
}

type MutableSetLike[T any] interface {
	SetLike[T]
	Add(T) error
	Del(T) error
	DelKey(string) error
	Reset(SetLike[T])
}

//  __     __    _              ____       _
//  \ \   / /_ _| |_   _  ___  / ___|  ___| |_ ___
//   \ \ / / _` | | | | |/ _ \ \___ \ / _ \ __/ __|
//    \ V / (_| | | |_| |  __/  ___) |  __/ |_\__ \
//     \_/ \__,_|_|\__,_|\___| |____/ \___|\__|___/
//

type Adder[E any] interface {
	Add(E) error
}

type Equaler[T any] interface {
	Equals(*T) bool
}

type Eacher[E any] interface {
	Each(FuncIter[E]) error
}

type EachPtrer[E any] interface {
	EachPtr(FuncIter[*E]) error
}

type StringAdder interface {
	AddString(string) error
}

type ValueSetLike[T flag.Value] interface {
	Strings() []string
	SetLike[T]
}

type MutableValueSetLike[T flag.Value] interface {
	MutableSetLike[T]
}
