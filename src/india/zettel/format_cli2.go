package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
)

// [kopf/schwanz@sha !typ]
func MakeCliFormatNamed(
	hf format.FormatWriterFunc[hinweis.Hinweis],
	sf format.FormatWriterFunc[sha.Sha],
	zf format.FormatWriterFunc[Zettel],
) format.FormatWriterFunc[Named] {
	return func(w io.Writer, z *Named) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			format.MakeWriter(hf, &z.Kennung),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &z.Stored.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter(zf, &z.Stored.Objekte),
			format.MakeFormatString("]"),
		)
	}
}
