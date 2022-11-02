package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

func MakeSerializedFormatWriter(
	f Format,
	out io.Writer,
	arf AkteReaderFactory,
	k konfig.Konfig,
) collections.WriterFunc[*Zettel] {
	wf := func(z *Zettel) (err error) {
		//use konfig
		includeAkte := z.Typ.String() == "md"

		if typKonfig, ok := k.Typen[z.Typ.String()]; ok {
			includeAkte = typKonfig.InlineAkte
		}

		ctx := FormatContextWrite{
			Out:               out,
			AkteReaderFactory: arf,
			Zettel:            *z,
			IncludeAkte:       includeAkte,
		}

		if _, err = f.WriteTo(ctx); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	return collections.MakeSyncSerializer(wf).Do
}
