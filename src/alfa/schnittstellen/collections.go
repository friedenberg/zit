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

type Lenner interface {
	Len() int
}

type ContainsKeyer interface {
	ContainsKey(string) bool
}

type SetLike[T any] interface {
	Lenner
	ContainsKeyer

	Equatable[Set[T]]
	Key(T) string
	Get(string) (T, bool)
	Any() T
	Contains(T) bool
	Each(FuncIter[T]) error
	EachPtr(FuncIter[*T]) error
	EachKey(FuncIterKey) error
	Elements() []T
}

type MutableSetLike[T any] interface {
	SetLike[T]
	Add(T) error
	Del(T) error
	DelKey(string) error
}

type Set[T any] interface {
	SetLike[T]
	ImmutableCloner[Set[T]]
	MutableCloner[MutableSet[T]]
}

type MutableSet[T any] interface {
	Set[T]
	MutableSetLike[T]
	Resetter2
}

type TridexLike interface {
	Lenner
	ContainsAbbreviation(string) bool
	ContainsExpansion(string) bool
	Abbreviate(string) string
	Expand(string) string
}

type MutableTridexLike interface {
	TridexLike
	Add(string)
	Remove(string)
}

type Tridex interface {
	TridexLike
	// ImmutableCloner[TridexLike]
	MutableCloner[MutableTridex]
}

type MutableTridex interface {
	Tridex
	Add(string)
	Remove(string)
}

// type Tridex interface {
//   TridexLike
// 	// ImmutableCloner[Tridex]
// 	MutableCloner[MutableTridex]
// }

// type MutableTridex interface {
// 	MutableTridexLike
// 	MutableCloner[Tridex]
// }

type Poolable[T any] interface{}

type PoolablePtr[T any] interface {
	Ptr[T]
	Resetable[T]
}

type Pool[T Poolable[T], TPtr PoolablePtr[T]] interface {
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
