package interfaces

import (
	"iter"

	"golang.org/x/exp/constraints"
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

type StringAdder interface {
	AddString(string) error
}

type Iterable[T any] interface {
	Any() T // TODO-P2 remove in favor of collection method
	// Each(FuncIter[T]) error // TODO remove in favor of iter.Seq
	All() iter.Seq[T]
	Lenner
}

type KeyedIterable[K constraints.Ordered, T any] interface {
	Iterable[T]
	AllPairs() iter.Seq2[K, T]
}

type IterablePtr[T any, TPtr Ptr[T]] interface {
	EachPtr(FuncIter[TPtr]) error // TODO remove in favor of iter.Seq
	AllPtr() iter.Seq[TPtr]
}
