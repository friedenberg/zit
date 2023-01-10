package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
)

type LogWriter[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] struct {
	New, Updated, Unchanged, Archived collections.WriterFunc[*Transacted[T, T1, T2, T3, T4, T5]]
}

type Store[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] interface {
	errors.Flusher
	AkteTextSaver[T, T1]
	TransactedInflator[T, T1, T2, T3, T4, T5]
	ReadOne(T3) (*Transacted[T, T1, T2, T3, T4, T5], error)
	ReadAllSchwanzen(collections.WriterFunc[*Transacted[T, T1, T2, T3, T4, T5]]) error
	ReadAll(collections.WriterFunc[*Transacted[T, T1, T2, T3, T4, T5]]) error
	SetLogWriter(LogWriter[T, T1, T2, T3, T4, T5])
}

type StoreWithCreateOrUpdate[
	T gattung.Objekte[T],
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
	T4 gattung.Verzeichnisse[T],
	T5 gattung.VerzeichnissePtr[T4, T],
] interface {
	CreateOrUpdate(T1, T3) (*Transacted[T, T1, T2, T3, T4, T5], error)
}
