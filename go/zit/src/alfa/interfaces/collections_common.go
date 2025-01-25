package interfaces

import (
	"iter"
)

type (
	FuncIter[T any]              func(T) error
	FuncIterIO[T any]            func(T) (int64, error)
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

type Iterable[T any] interface {
	Any() T
	All() iter.Seq[T]
}

type IterablePtr[T any, TPtr Ptr[T]] interface {
	EachPtr(FuncIter[TPtr]) error // TODO remove in favor of iter.Seq
	AllPtr() iter.Seq[TPtr]
}
