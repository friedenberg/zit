package zettel_transacted

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
)

// (created|updated|archived) [kopf/schwanz@sha !typ]
// TODO add archived state
func MakeCliFormat(
	znf collections.WriterFuncFormat[zettel_named.Zettel],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		verb := ""

		switch {
		case z.IsNew():
			verb = "created"

		default:
			verb = "updated"
		}

		return collections.WriteFormats(
			w,
			collections.MakeWriterLiteral("("),
			collections.MakeWriterLiteral(verb),
			collections.MakeWriterLiteral(")"),
			collections.MakeWriterLiteral(" "),
			collections.MakeWriterFormatFunc(znf, &z.Named),
		)
	}
}
