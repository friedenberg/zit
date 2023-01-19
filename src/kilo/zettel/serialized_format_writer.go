package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

func MakeSerializedFormatWriter(
	f ObjekteFormatter,
	out io.Writer,
	arf schnittstellen.AkteReaderFactory,
	k konfig.Compiled,
) collections.WriterFunc[*Objekte] {
	errors.TodoP3("rename to MakeSingleplexedFormatWriter")

	wf := func(z *Objekte) (err error) {
		c := ObjekteFormatterContext{
			Zettel:      *z,
			IncludeAkte: k.IsInlineTyp(z.Typ),
		}

		if _, err = f.Format(out, &c); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	return collections.MakeSyncSerializer(wf)
}
