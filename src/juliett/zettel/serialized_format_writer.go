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
) schnittstellen.FuncIter[*Transacted] {
	errors.TodoP3("rename to MakeSingleplexedFormatWriter")

	wf := func(z *Transacted) (err error) {
		c := ObjekteFormatterContext{
			Zettel:      z.Objekte,
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
