package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/transacted"
)

// [kopf/schwanz@sha !typ]
func MakeCliFormatTransacted(
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[transacted.Zettel] {
	return func(w io.Writer, z transacted.Zettel) (n int64, err error) {
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
	ztf schnittstellen.FuncWriterFormat[transacted.Zettel],
) schnittstellen.FuncWriterFormat[transacted.Zettel] {
	return func(w io.Writer, z transacted.Zettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeWriter(ztf, z),
		)
	}
}

type LogWriter = objekte_store.LogWriter[*transacted.Zettel]
