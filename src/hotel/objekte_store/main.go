package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
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

type ExternalReader[
	E any,
	V any,
] interface {
	ReadOneExternal(E) (V, error)
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
	K any,
	V any,
] interface {
	TransactedReader[K, V]
	Query(kennung.Matcher, schnittstellen.FuncIter[V]) error
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

type CheckedOutUpdater[
	CO objekte.CheckedOutLike,
	T objekte.TransactedLike,
] interface {
	UpdateCheckedOut(CO) (T, error)
}

type CreateOrUpdater[
	O any,
	K any,
	V any,
	CO any,
] interface {
	CreateOrUpdate(O, K) (V, error)
	CreateOrUpdateCheckedOut(CO) (V, error)
}
