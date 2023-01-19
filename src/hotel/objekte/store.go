package objekte

import (
	"github.com/friedenberg/zit/src/charlie/collections"
)

type LogWriter[
	T any,
] struct {
	New, Updated, Unchanged, Archived collections.WriterFunc[T]
}

type TransactedLogger[
	T any,
] interface {
	SetLogWriter(LogWriter[T])
}

type Reader[
	K any,
	V any,
] interface {
	ReadOne(K) (V, error)
	ReadAllSchwanzen(collections.WriterFunc[V]) error
}

type TransactedReader[
	K any,
	V any,
] interface {
	Reader[K, V]
	ReadAll(collections.WriterFunc[V]) error
}

type CreateOrUpdater[
	O any,
	K any,
	V any,
] interface {
	CreateOrUpdate(O, K) (V, error)
}

type Updater[
	O any,
	V any,
] interface {
	Update(O) (V, error)
}
