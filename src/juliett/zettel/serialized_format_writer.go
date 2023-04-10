package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/india/konfig"
)

func MakeSerializedFormatWriter(
	f ObjekteFormatter,
	out io.Writer,
	arf schnittstellen.AkteReaderFactory,
	k konfig.Compiled,
) schnittstellen.FuncIter[*Objekte] {
	errors.TodoP3("rename to MakeSingleplexedFormatWriter")

	wf := func(z *Objekte) (err error) {
		c := ObjekteFormatterContext{
			Zettel:      *z,
			IncludeAkte: k.IsInlineTyp(z.GetTyp()),
		}

		if _, err = f.Format(out, &c); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	return collections.MakeSyncSerializer(wf)
}
