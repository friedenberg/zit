package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/india/zettel"
)

// (new|unchanged|updated|archived) [kopf/schwanz@sha !typ]
func MakeCliFormat(
	znf format.FormatWriterFunc[zettel.Named],
	verb string,
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(verb),
			format.MakeWriter(znf, &z.Named),
		)
	}
}
