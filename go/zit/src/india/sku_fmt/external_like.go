package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type (
	ReaderExternalLike = catgut.StringFormatReader[sku.ExternalLike]
	WriterExternalLike = catgut.StringFormatWriter[sku.ExternalLike]

	ExternalLike       struct {
		ReaderExternalLike
		WriterExternalLike
	}
)
