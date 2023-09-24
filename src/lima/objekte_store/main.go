package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
)

type TransactedLogger interface {
	SetLogWriter(LogWriter[sku.SkuLikePtr])
}

type LastReader interface {
	ReadLast() (*sku.Transacted, error)
}

type OneReader interface {
	ReadOne(schnittstellen.StringerGattungGetter) (*sku.Transacted, error)
}

type ExternalReader[
	E any,
	T any,
] interface {
	ReadOneExternal(E, T) (*sku.External, error)
}

type AllReader interface {
	ReadAll(schnittstellen.FuncIter[*sku.Transacted]) error
}

type SchwanzenReader interface {
	ReadAllSchwanzen(schnittstellen.FuncIter[*sku.Transacted]) error
}

type Reader interface {
	OneReader
	SchwanzenReader
}

type TransactedReader interface {
	Reader
	AllReader
}

type Querier interface {
	TransactedReader
	Query(matcher.MatcherSigil, schnittstellen.FuncIter[*sku.Transacted]) error
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
	CO *sku.CheckedOut,
	T sku.SkuLike,
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
		metadatei.Getter,
		K,
		schnittstellen.ShaLike,
	) (V, error)
	CreateOrUpdate(metadatei.Getter, K) (V, error)
	CreateOrUpdateCheckedOut(CO) (V, error)
}
