package collections

import (
	"flag"
	"fmt"

	"github.com/friedenberg/zit/src/charlie/gattung"
)

//   _____
//  |  ___|   _ _ __   ___ ___
//  | |_ | | | | '_ \ / __/ __|
//  |  _|| |_| | | | | (__\__ \
//  |_|   \__,_|_| |_|\___|___/
//

type WriterFunc[T any] func(T) error

// type WriterFuncFuncPooled[T any] func(PoolLike[T]) WriterFunc[*T]
type WriterFuncWithKey[T any] func(string, T) error
type WriterFuncKey func(string) error

//   ____             _
//  |  _ \ ___   ___ | |___
//  | |_) / _ \ / _ \| / __|
//  |  __/ (_) | (_) | \__ \
//  |_|   \___/ \___/|_|___/
//

type PoolLike[T any] interface {
	Get() *T
	Put(i *T) (err error)
}

type Pool2Like[T gattung.Element, TPtr gattung.ElementPtr[T]] interface {
	Get() TPtr
	Put(i TPtr) (err error)
}

//   ____       _
//  / ___|  ___| |_ ___
//  \___ \ / _ \ __/ __|
//   ___) |  __/ |_\__ \
//  |____/ \___|\__|___/
//

type SetLike[T any] interface {
	Len() int
	Key(T) string
	Get(string) (T, bool)
	ContainsKey(string) bool
	Contains(T) bool
	Each(WriterFunc[T]) error
	EachPtr(WriterFunc[*T]) error
	EachKey(WriterFuncKey) error
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

type ValueSetElement interface {
	gattung.Element
	fmt.Stringer
}

type ValueSetElementPtr[E gattung.Element] interface {
	gattung.ElementPtr[E]
	Setter
}

type FuncSetString = gattung.FuncSetString

type Setter interface {
	Set(string) error
}

type SetterPtr[T any] interface {
	gattung.ElementPtr[T]
	Setter
}

type Adder[E any] interface {
	Add(E) error
}

type Eacher[E any] interface {
	Each(WriterFunc[E]) error
}

type EachPtrer[E any] interface {
	EachPtr(WriterFunc[*E]) error
}

type StringAdder interface {
	AddString(string) error
}

type ValueSetLike[T flag.Value] interface {
	SetLike[T]
}

type MutableValueSetLike[T flag.Value] interface {
	MutableSetLike[T]
}
