package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type TransactedLogger[
	T any,
] interface {
	SetLogWriter(LogWriter[objekte.TransactedLikePtr])
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
	T any,
	V any,
] interface {
	ReadOneExternal(E, T) (V, error)
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
	Query(kennung.MatcherSigil, schnittstellen.FuncIter[V]) error
}

type Creator[
	V any,
] interface {
	Create(metadatei.Getter) (V, error)
}

type Updater[
	K any,
	V any,
] interface {
	Update(metadatei.Getter, K) (V, error)
}

type CheckedOutUpdater[
	CO objekte.CheckedOutLike,
	T objekte.TransactedLike,
] interface {
	UpdateCheckedOut(CO) (T, error)
}

type UpdaterManyMetadatei interface {
	UpdateManyMetadatei(
		schnittstellen.SetLike[sku.SkuLike],
	) error
}

type CreateOrUpdater[
	O any,
	K any,
	V any,
	CO any,
] interface {
	CreateOrUpdateAkte(
		O,
		metadatei.Getter,
		K,
		schnittstellen.ShaLike,
	) (V, error)
	CreateOrUpdate(O, metadatei.Getter, K) (V, error)
	CreateOrUpdateCheckedOut(CO) (V, error)
}
