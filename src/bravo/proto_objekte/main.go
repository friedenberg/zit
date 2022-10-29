package proto_objekte

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

type SetLike[T any] interface {
	Len() int
	Contains(T) bool
	WriterContainer() WriterFunc[T]
	Each(WriterFunc[T]) error
}

type MutableSetLike[T any] interface {
	SetLike[T]
	WriterAdder() WriterFunc[T]
	WriterRemover() WriterFunc[T]
	Reset(SetLike[T])
}
