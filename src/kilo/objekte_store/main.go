package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
)

type (
	TransactedLogger interface {
		SetLogWriter(LogWriter)
	}

	LastReader interface {
		ReadLast() (*sku.Transacted, error)
	}

	OneReader interface {
		ReadOne(schnittstellen.StringerGattungGetter) (*sku.Transacted, error)
	}

	ExternalReader[
		E any,
		T any,
	] interface {
		ReadOneExternal(E, T) (*sku.External, error)
	}

	AllReader interface {
		ReadAll(schnittstellen.FuncIter[*sku.Transacted]) error
	}

	SchwanzenReader interface {
		ReadAllSchwanzen(schnittstellen.FuncIter[*sku.Transacted]) error
	}

	Reader interface {
		OneReader
		SchwanzenReader
	}

	TransactedReader interface {
		Reader
		AllReader
	}

	Querier interface {
		TransactedReader
		Query(
			matcher.MatcherSigil,
			schnittstellen.FuncIter[*sku.Transacted],
		) error
	}

	Creator[
		V any,
	] interface {
		Create(metadatei.Getter) (V, error)
	}

	Updater[
		K any,
		V any,
	] interface {
		Update(metadatei.Getter, K) (V, error)
	}

	CheckedOutUpdater interface {
		UpdateCheckedOut(*sku.CheckedOut) (*sku.Transacted, error)
	}

	CreateOrUpdater interface {
		CreateOrUpdateAkte(
			metadatei.Getter,
			kennung.Kennung,
			schnittstellen.ShaLike,
		) (*sku.Transacted, error)
		CreateOrUpdate(
			metadatei.Getter,
			kennung.Kennung,
		) (*sku.Transacted, error)
		CreateOrUpdateCheckedOut(*sku.CheckedOut) (*sku.Transacted, error)
	}
)
