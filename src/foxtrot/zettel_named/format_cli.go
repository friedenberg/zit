package zettel_named

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/zettel"
)

// [kopf/schwanz@sha !typ]
func MakeCliFormat(
	hf collections.WriterFuncFormat[hinweis.Hinweis],
	sf collections.WriterFuncFormat[sha.Sha],
	zf collections.WriterFuncFormat[zettel.Zettel],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return collections.WriteFormats(
			w,
			collections.MakeWriterLiteral("["),
			collections.MakeWriterFormatFunc(hf, &z.Hinweis),
			collections.MakeWriterLiteral("@"),
			collections.MakeWriterFormatFunc(sf, &z.Stored.Sha),
			collections.MakeWriterLiteral(" "),
			collections.MakeWriterFormatFunc(zf, &z.Stored.Zettel),
			collections.MakeWriterLiteral("]"),
		)
	}
}
