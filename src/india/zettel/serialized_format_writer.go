package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/hotel/konfig_compiled"
	"github.com/friedenberg/zit/src/hotel/typ"
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
			Out:               out,
			AkteReaderFactory: arf,
			Zettel:            *z,
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
