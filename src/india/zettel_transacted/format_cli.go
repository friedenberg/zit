package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/format"
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
			verb = "created"

		default:
			verb = "updated"
		}

		return format.Write(
			w,
			format.MakeFormatString("(%s) ", verb),
			format.MakeWriter(znf, &z.Named),
		)
	}
}
