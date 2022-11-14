package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
)

// (created|updated|archived) [kopf/schwanz@sha !typ]
// TODO add archived state
func MakeCliFormat(
	znf format.FormatWriterFunc[zettel_named.Zettel],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		verb := ""

		switch {
		case z.IsNew():
      verb = format.StringNew

		default:
      verb = format.StringUpdated
		}

		return format.Write(
			w,
      format.MakeFormatStringRightAlignedParen(verb),
			format.MakeWriter(znf, &z.Named),
		)
	}
}
