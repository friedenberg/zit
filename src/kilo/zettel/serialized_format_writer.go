package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

func MakeSerializedFormatWriter(
	f Format,
	out io.Writer,
	arf gattung.AkteReaderFactory,
	k konfig_compiled.Compiled,
) collections.WriterFunc[*Objekte] {
	wf := func(z *Objekte) (err error) {
		isInline := typ.IsInlineAkte(z.Typ, k)

		ctx := FormatContextWrite{
			Out:    out,
			Zettel: *z,
			//TODO this seems inverted for some reason
			IncludeAkte: isInline,
		}

		if _, err = f.WriteTo(ctx); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	return collections.MakeSyncSerializer(wf)
}
