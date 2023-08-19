package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

// [kopf/schwanz@sha !typ]
func MakeCliFormatTransacted(
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[sku.TransactedZettel] {
	return func(w io.Writer, z sku.TransactedZettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			format.MakeWriter(hf, z.Kennung),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.GetAkteSha()),
			format.MakeFormatString(" "),
			format.MakeWriter[metadatei.GetterPtr](mf, &z),
			format.MakeFormatString("]"),
		)
	}
}

// (new|unchanged|updated|archived) [kopf/schwanz@sha !typ]
func MakeCliFormatTransactedDelta(
	ztf schnittstellen.FuncWriterFormat[sku.TransactedZettel],
) schnittstellen.FuncWriterFormat[sku.TransactedZettel] {
	return func(w io.Writer, z sku.TransactedZettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeWriter(ztf, z),
		)
	}
}

type LogWriter = objekte_store.LogWriter[*sku.TransactedZettel]
