package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/golf/typ"
)

func MakeSerializedFormatWriter(
	f Format,
	out io.Writer,
	arf gattung.AkteReaderFactory,
	k konfig.Konfig,
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
