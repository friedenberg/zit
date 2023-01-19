package collections

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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

type Pool2Like[T any, TPtr schnittstellen.Ptr[T]] interface {
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

type Adder[E any] interface {
	Add(E) error
}

type Lessor[T any] interface {
	Less(T) bool
}

type Equaler[T any] interface {
	Equals(*T) bool
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
