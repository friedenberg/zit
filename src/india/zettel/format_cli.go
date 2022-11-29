package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/delta/kennung"
)

// !typ "bez"
func MakeCliFormat(
	bf format.FormatWriterFunc[bezeichnung.Bezeichnung],
	ef format.FormatWriterFunc[kennung.EtikettSet],
	tf format.FormatWriterFunc[kennung.Typ],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
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
	znf format.FormatWriterFunc[Named],
	verb string,
) format.FormatWriterFunc[Transacted] {
	return func(w io.Writer, z *Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(verb),
			format.MakeWriter(znf, &z.Named),
		)
	}
}
