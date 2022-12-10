package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
)

// !typ "bez"
func MakeCliFormat(
	bf format.FormatWriterFunc[bezeichnung.Bezeichnung],
	ef format.FormatWriterFunc[kennung.EtikettSet],
	tf format.FormatWriterFunc[kennung.Typ],
) format.FormatWriterFunc[Objekte] {
	return func(w io.Writer, z *Objekte) (n int64, err error) {
		var lastWriter format.WriterFunc

		if z.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(ef, &z.Etiketten)
		} else {
			lastWriter = format.MakeWriter(bf, &z.Bezeichnung)
		}

		return format.Write(
			w,
			format.MakeWriter(tf, &z.Typ),
			format.MakeFormatString(" "),
			lastWriter,
		)
	}
}

// (new|unchanged|updated|archived) [kopf/schwanz@sha !typ]
func MakeCliFormatTransacted(
	hf format.FormatWriterFunc[hinweis.Hinweis],
	sf format.FormatWriterFunc[sha.Sha],
	zf format.FormatWriterFunc[Objekte],
	verb string,
) format.FormatWriterFunc[Transacted] {
	return func(w io.Writer, z *Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(verb),
			format.MakeFormatString("["),
			format.MakeWriter(hf, z.Kennung()),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &z.Sku.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter[Objekte](zf, &z.Objekte),
			format.MakeFormatString("]"),
		)
	}
}
