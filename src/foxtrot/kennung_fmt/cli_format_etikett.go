package kennung_fmt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type etikettenCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeEtikettenCliFormat() *etikettenCliFormat {
	return &etikettenCliFormat{
		stringFormatWriter: string_format_writer.MakeString[string](),
	}
}

func (f *etikettenCliFormat) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	k kennung.EtikettSet,
) (n int64, err error) {
	if k.Len() == 0 {
		return
	}

	v := iter.StringDelimiterSeparated[kennung.Etikett](" ", k)

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
