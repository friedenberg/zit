package objekte

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

func MakeCliFormatTransactedLike(
	hf schnittstellen.FuncWriterFormat[kennung.IdLike],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	mf schnittstellen.FuncWriterFormat[metadatei.Metadatei],
) schnittstellen.FuncWriterFormat[TransactedLike] {
	return func(w io.Writer, z TransactedLike) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			format.MakeWriter(hf, z.GetKennung()),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.GetObjekteSha()),
			format.MakeFormatString(" "),
			format.MakeWriter[metadatei.Metadatei](mf, z.GetMetadatei()),
			format.MakeFormatString("]"),
		)
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
