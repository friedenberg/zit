package collections

import (
	"flag"
	"fmt"
)

type ProtoObjekte interface {
	fmt.Stringer
}

type ProtoObjektePointer interface {
	flag.Value
}

type WriterFunc[T any] func(T) error
type WriterFuncWithKey[T any] func(string, T) error
type WriterFuncKey func(string) error

type SetLike[T any] interface {
	Len() int
	Key(T) string
	Get(string) (T, bool)
	ContainsKey(string) bool
	Contains(T) bool
	WriterContainer() WriterFunc[T]
	Each(WriterFunc[T]) error
	EachKey(WriterFuncKey) error
	Elements() []T
}

type MutableSetLike[T any] interface {
	SetLike[T]
	Add(T) error
	Del(T) error
	DelKey(string) error
	Reset(SetLike[T])
}

type ValueSetLike[T flag.Value] interface {
	SetLike[T]
}
