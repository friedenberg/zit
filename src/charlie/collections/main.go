package collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

//   _____
//  |  ___|   _ _ __   ___ ___
//  | |_ | | | | '_ \ / __/ __|
//  |  _|| |_| | | | | (__\__ \
//  |_|   \__,_|_| |_|\___|___/
//

// type WriterFuncFuncPooled[T any] func(PoolLike[T]) WriterFunc[*T]
type (
	WriterFuncWithKey[T any] func(string, T) error
	WriterFuncKey            func(string) error
)

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
