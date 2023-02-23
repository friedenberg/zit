package schnittstellen

import (
	"flag"
)

type (
	FuncIter[T any]              func(T) error
	FuncTransform[T any, T1 any] func(T) (T1, error)
	FuncIterKey                  func(string) error
	FuncIterWithKey[T any]       func(string, T) error
)

type Set[T any] interface {
	Equatable[Set[T]]
	Len() int
	Key(T) string
	Get(string) (T, bool)
	Any() T
	ContainsKey(string) bool
	Contains(T) bool
	Each(FuncIter[T]) error
	EachPtr(FuncIter[*T]) error
	EachKey(FuncIterKey) error
	ImmutableCloner[Set[T]]
	MutableCloner[MutableSet[T]]
	Elements() []T
}

type MutableSet[T any] interface {
	Set[T]
	Add(T) error
	Del(T) error
	DelKey(string) error
	Resetter2
	// ResetWither[MutableSet[T]]
}

type Pool[T any, TPtr Ptr[T]] interface {
	Get() TPtr
	Put(i TPtr) (err error)
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

type ValueSet[T flag.Value] interface {
	Strings() []string
	Set[T]
}

type MutableValueSet[T flag.Value] interface {
	MutableSet[T]
}
