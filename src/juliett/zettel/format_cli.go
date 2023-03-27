package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

// !typ "bez"
func MakeCliFormat(
	bf schnittstellen.FuncWriterFormat[bezeichnung.Bezeichnung],
	ef schnittstellen.FuncWriterFormat[schnittstellen.SetLike[kennung.Etikett]],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
) schnittstellen.FuncWriterFormat[Objekte] {
	return func(w io.Writer, z Objekte) (n int64, err error) {
		var lastWriter schnittstellen.FuncWriter

		if z.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(ef, schnittstellen.SetLike[kennung.Etikett](z.Etiketten))
		} else {
			lastWriter = format.MakeWriter(bf, z.Bezeichnung)
		}

		return format.Write(
			w,
			format.MakeWriter(tf, z.Typ),
			format.MakeFormatString(" "),
			lastWriter,
		)
	}
}

// [kopf/schwanz@sha !typ]
func MakeCliFormatTransacted(
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	zf schnittstellen.FuncWriterFormat[Objekte],
) schnittstellen.FuncWriterFormat[Transacted] {
	return func(w io.Writer, z Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			format.MakeWriter(hf, *z.Kennung()),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.GetObjekteSha()),
			format.MakeFormatString(" "),
			format.MakeWriter[Objekte](zf, z.Objekte),
			format.MakeFormatString("]"),
		)
	}
}

// (new|unchanged|updated|archived) [kopf/schwanz@sha !typ]
func MakeCliFormatTransactedDelta(
	ztf schnittstellen.FuncWriterFormat[Transacted],
) schnittstellen.FuncWriterFormat[Transacted] {
	return func(w io.Writer, z Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeWriter(ztf, z),
		)
	}
}

type LogWriter = objekte_store.LogWriter[*Transacted]
