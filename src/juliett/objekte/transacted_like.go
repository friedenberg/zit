package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
)

type (
	FuncReaderTransactedLikePtr  func(schnittstellen.FuncIter[*sku.Transacted]) error
	FuncQuerierTransactedLikePtr func(matcher.MatcherSigil, schnittstellen.FuncIter[*sku.Transacted]) error
)
