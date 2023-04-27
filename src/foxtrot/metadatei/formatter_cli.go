package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
)

// !typ "bez"
func MakeCliFormatExcludeTyp(
	bf schnittstellen.FuncWriterFormat[bezeichnung.Bezeichnung],
	ef schnittstellen.FuncWriterFormat[schnittstellen.SetLike[kennung.Etikett]],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
) schnittstellen.FuncWriterFormat[GetterPtr] {
	return func(w io.Writer, z GetterPtr) (n int64, err error) {
		m := z.GetMetadateiPtr()

		var lastWriter schnittstellen.FuncWriter

		if m.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(
				ef,
				schnittstellen.SetLike[kennung.Etikett](m.Etiketten),
			)
		} else if !m.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(bf, m.Bezeichnung)
		} else {
			return
		}

		return format.Write(
			w,
			lastWriter,
		)
	}
}

func MakeCliFormatIncludeTyp(
	bf schnittstellen.FuncWriterFormat[bezeichnung.Bezeichnung],
	ef schnittstellen.FuncWriterFormat[schnittstellen.SetLike[kennung.Etikett]],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
) schnittstellen.FuncWriterFormat[GetterPtr] {
	return func(w io.Writer, z GetterPtr) (n int64, err error) {
		m := z.GetMetadateiPtr()

		var lastWriter schnittstellen.FuncWriter

		if m.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(
				ef,
				schnittstellen.SetLike[kennung.Etikett](m.Etiketten),
			)
		} else if !m.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(bf, m.Bezeichnung)
		} else {
			return
		}

		return format.Write(
			w,
			format.MakeWriter(tf, m.GetTyp()),
			format.MakeFormatString(" "),
			lastWriter,
		)
	}
}
