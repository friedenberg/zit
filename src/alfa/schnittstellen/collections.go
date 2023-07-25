package schnittstellen

import (
	"flag"
)

type (
	FuncIter[T any]              func(T) error
	FuncIterIO[T any]            func(T) (int64, error)
	FuncTransform[T any, T1 any] func(T) (T1, error)
	FuncIterKey                  func(string) error
	FuncIterWithKey[T any]       func(string, T) error
)

type Delta[T ValueLike] interface {
	GetAdded() Set[T]
	GetRemoved() Set[T]
}

type Lenner interface {
	Len() int
}

type Iterable[T any] interface {
	Any() T // TODO-P2 remove in favor of collection method?
	Each(FuncIter[T]) error
	// EachPtr(FuncIter[*T]) error
}

type IterablePtr[T any, TPtr Ptr[T]] interface {
	EachPtr(FuncIter[TPtr]) error
}

type Collection[T any] interface {
	Lenner
	Iterable[T]
}

type ContainsKeyer interface {
	ContainsKey(string) bool
}

type SetLike[T any] interface {
	Collection[T]
	ContainsKeyer

	EqualsSetLike(SetLike[T]) bool
	Key(T) string
	Get(string) (T, bool)
	Contains(T) bool
	EachKey(FuncIterKey) error
	Elements() []T // TODO-P2 remove in favor of collection method

	CloneSetLike() SetLike[T]
	CloneMutableSetLike() MutableSetLike[T]
}

type MutableSetLike[T any] interface {
	SetLike[T]
	Adder[T]
	AdderCustom[T]
	Del(T) error
	DelKey(string) error
	Resetter
}

type Set[T any] interface {
	SetLike[T]
}

type MutableSet[T any] interface {
	Set[T]
	MutableSetLike[T]
}

type TridexLike interface {
	Lenner
	EachString(FuncIter[string]) error
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
}

type MutableTridex interface {
	Tridex
	Add(string)
	Remove(string)
}

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

type AdderCustom[E any] interface {
	AddCustomKey(E, func(E) string) error
}

// TODO-P1 remove
type Equaler[T any] interface {
	Equals(*T) bool
}

// TODO-P1 remove
type Eacher[E any] interface {
	Each(FuncIter[E]) error
}

// TODO-P1 remove
type EachPtrer[E any] interface {
	EachPtr(FuncIter[*E]) error
}

type StringAdder interface {
	AddString(string) error
}

type ValueSet[T flag.Value, TPtr ValuePtr[T]] interface {
	Strings() []string
	Set[T]
}

type MutableValueSet[T flag.Value, TPtr ValuePtr[T]] interface {
	MutableSet[T]
}
