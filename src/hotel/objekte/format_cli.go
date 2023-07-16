package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

func MakeCliFormatTransactedLikePtr(
	hf schnittstellen.FuncWriterFormat[kennung.Kennung],
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[TransactedLikePtr] {
	return func(w io.Writer, z TransactedLikePtr) (n int64, err error) {
		if z.GetMetadateiPtr().UserInputIsEmpty() {
			return format.Write(
				w,
				format.MakeFormatString("["),
				format.MakeWriter(hf, z.GetKennung()),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetAkteSha()),
				format.MakeFormatString("]"),
			)
		} else {
			return format.Write(
				w,
				format.MakeFormatString("["),
				format.MakeWriter(hf, z.GetKennung()),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetAkteSha()),
				format.MakeFormatString(" "),
				format.MakeWriter[metadatei.GetterPtr](mf, z),
				format.MakeFormatString("]"),
			)
		}
	}
}

// (new|unchanged|updated|archived) [kopf/schwanz@sha !typ]
func MakeCliFormatTransactedLikeDelta(
	ztf schnittstellen.FuncWriterFormat[TransactedLike],
) schnittstellen.FuncWriterFormat[TransactedLike] {
	return func(w io.Writer, z TransactedLike) (n int64, err error) {
		return format.Write(
			w,
			format.MakeWriter(ztf, z),
		)
	}
}
