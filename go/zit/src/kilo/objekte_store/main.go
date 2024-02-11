package objekte_store

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/matcher"
)

type (
	LastReader interface {
		ReadLast() (*sku.Transacted, error)
	}

	OneReader interface {
		ReadOne(schnittstellen.StringerGattungGetter) (*sku.Transacted, error)
	}

	ExternalReader interface {
		ReadOneExternal(
			*sku.ExternalMaybe,
			*sku.Transacted,
		) (*sku.External, error)
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

	AkteTextSaver[
		T schnittstellen.Akte[T],
		T1 schnittstellen.AktePtr[T],
	] interface {
		SaveAkteText(T1) (schnittstellen.ShaLike, int64, error)
	}
)
