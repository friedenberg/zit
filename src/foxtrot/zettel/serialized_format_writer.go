package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/konfig"
)

func MakeSerializedFormatWriter(
	f Format,
	out io.Writer,
	arf AkteReaderFactory,
	k konfig.Konfig,
) collections.WriterFunc[*Zettel] {
	wf := func(z *Zettel) (err error) {
		isInline := z.Typ.IsInline(k)

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
