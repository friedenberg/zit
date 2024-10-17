package interfaces

import "iter"

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

type StringKeyer[T any] interface {
	GetKey(T) string
}

type StringKeyerPtr[T any, TPtr Ptr[T]] interface {
	StringKeyer[T]
	GetKeyPtr(TPtr) string
}

type ContainsKeyer interface {
	ContainsKey(string) bool
}

type StringAdder interface {
	AddString(string) error
}

type Iterable[T any] interface {
	Any() T                 // TODO-P2 remove in favor of collection method
	Each(FuncIter[T]) error // TODO remove in favor of iter.Seq
	All() iter.Seq[T]
	Lenner
}

type IterablePtr[T any, TPtr Ptr[T]] interface {
	EachPtr(FuncIter[TPtr]) error // TODO remove in favor of iter.Seq
	AllPtr() iter.Seq[TPtr]
}
