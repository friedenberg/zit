package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type TransactedLogger[
	T any,
] interface {
	SetLogWriter(LogWriter[T])
}

type LastReader[
	V any,
] interface {
	ReadLast() (V, error)
}

type OneReader[
	K any,
	V any,
] interface {
	ReadOne(K) (V, error)
}

type AllReader[
	V any,
] interface {
	ReadAll(schnittstellen.FuncIter[V]) error
}

type SchwanzenReader[
	V any,
] interface {
	ReadAllSchwanzen(schnittstellen.FuncIter[V]) error
}

type Reader[
	K any,
	V any,
] interface {
	OneReader[K, V]
	SchwanzenReader[V]
}

type TransactedReader[
	K any,
	V any,
] interface {
	Reader[K, V]
	AllReader[V]
}

type Querier[
	V any,
] interface {
	Query(kennung.Set, schnittstellen.FuncIter[V]) error
}

type CreateOrUpdater[
	O any,
	K any,
	V any,
] interface {
	CreateOrUpdate(O, K) (V, error)
}

type Creator[
	O any,
	V any,
] interface {
	Create(O) (V, error)
}

type Updater[
	O any,
	K any,
	V any,
] interface {
	Update(O, K) (V, error)
}
