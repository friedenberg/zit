package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
)

type TransactedLike interface {
	metadatei.Getter
	GetKennung() kennung.Kennung
	GetAkteSha() schnittstellen.Sha
	GetSku() sku.Sku
	GetSkuLike() sku.SkuLike
	kennung.Matchable
	sku.DataIdentityGetter
}

type TransactedLikePtr interface {
	TransactedLike
	metadatei.GetterPtr
	metadatei.Setter
	GetMetadateiWithKennung() metadatei.WithKennung
	GetKennungPtr() kennung.KennungPtr
	SetTai(kennung.Tai)
}

type StoredLikePtr interface {
	metadatei.Getter
	metadatei.Setter
	GetAkteSha() schnittstellen.Sha
	SetAkteSha(schnittstellen.Sha)
	SetObjekteSha(schnittstellen.Sha)
}

type (
	FuncReaderTransacted[T TransactedLike]       func(schnittstellen.FuncIter[T]) error
	FuncReaderTransactedPtr[T TransactedLikePtr] func(schnittstellen.FuncIter[T]) error
	FuncReaderTransactedLike                     func(schnittstellen.FuncIter[TransactedLike]) error
	FuncReaderTransactedLikePtr                  func(schnittstellen.FuncIter[TransactedLikePtr]) error
)

type (
	FuncQuerierTransacted[T TransactedLike]       func(kennung.Matcher, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedPtr[T TransactedLikePtr] func(kennung.Matcher, schnittstellen.FuncIter[T]) error
	FuncQuerierTransactedLike                     func(kennung.Matcher, schnittstellen.FuncIter[TransactedLike]) error
	FuncQuerierTransactedLikePtr                  func(kennung.Matcher, schnittstellen.FuncIter[TransactedLikePtr]) error
)

func MakeApplyQueryTransactedLikePtr[T TransactedLikePtr](
	fat FuncQuerierTransactedPtr[T],
) FuncQuerierTransactedLikePtr {
	return func(ids kennung.Matcher, fatl schnittstellen.FuncIter[TransactedLikePtr]) (err error) {
		return fat(
			ids,
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}

func MakeApplyTransactedLikePtr[T TransactedLikePtr](
	fat FuncReaderTransacted[T],
) FuncReaderTransactedLikePtr {
	return func(fatl schnittstellen.FuncIter[TransactedLikePtr]) (err error) {
		return fat(
			func(e T) (err error) {
				return fatl(e)
			},
		)
	}
}
