package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

// !typ "bez"
func MakeCliFormat(
	bf schnittstellen.FuncWriterFormat[bezeichnung.Bezeichnung],
	ef schnittstellen.FuncWriterFormat[schnittstellen.SetLike[kennung.Etikett]],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
) schnittstellen.FuncWriterFormat[Metadatei] {
	return func(w io.Writer, z Metadatei) (n int64, err error) {
		var lastWriter schnittstellen.FuncWriter

		if z.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(
				ef,
				schnittstellen.SetLike[kennung.Etikett](z.Etiketten),
			)
		} else {
			lastWriter = format.MakeWriter(bf, z.Bezeichnung)
		}

		return format.Write(
			w,
			format.MakeWriter(tf, z.GetTyp()),
			format.MakeFormatString(" "),
			lastWriter,
		)
	}
}
