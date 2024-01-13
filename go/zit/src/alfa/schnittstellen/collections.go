package schnittstellen

import (
	"flag"
)

type Delta[T any] interface {
	GetAdded() SetLike[T]
	GetRemoved() SetLike[T]
}

type Collection[T any] interface {
	Lenner
	Iterable[T]
}

type SetLike[T any] interface {
	Collection[T]
	ContainsKeyer

	Key(T) string
	Get(string) (T, bool)
	Contains(T) bool
	EachKey(FuncIterKey) error

	CloneSetLike() SetLike[T]
	CloneMutableSetLike() MutableSetLike[T]
}

type MutableSetLike[T any] interface {
	SetLike[T]
	Adder[T]
	Del(T) error
	DelKey(string) error
	Resetter
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
	// Resetable[T]
}

type PoolValue[T any] interface {
	Get() T
	Put(i T) (err error)
}

type Pool[T Poolable[T], TPtr PoolablePtr[T]] interface {
	PoolValue[TPtr]
	PutMany(...TPtr) error
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

type AdderPtr[E any, EPtr Ptr[E]] interface {
	AddPtr(EPtr) error
}

type AdderCustom[E any] interface {
	AddCustomKey(E, func(E) string) error
}

type ValueSet[T flag.Value, TPtr ValuePtr[T]] interface {
	Strings() []string
	SetLike[T]
}

type MutableValueSet[T flag.Value, TPtr ValuePtr[T]] interface {
	MutableSetLike[T]
}
