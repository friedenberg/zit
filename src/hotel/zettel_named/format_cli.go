package zettel_named

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

// [kopf/schwanz@sha !typ]
func MakeCliFormat(
	hf format.FormatWriterFunc[hinweis.Hinweis],
	sf format.FormatWriterFunc[sha.Sha],
	zf format.FormatWriterFunc[zettel.Zettel],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			format.MakeWriter(hf, &z.Hinweis),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &z.Stored.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter(zf, &z.Stored.Objekte),
			format.MakeFormatString("]"),
		)
	}
}
