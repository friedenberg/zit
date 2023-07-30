package kennung

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/string_writer_format"
)

type fdCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeFDCliFormat(
	co string_writer_format.ColorOptions,
	relativePathStringFormatWriter schnittstellen.StringFormatWriter[string],
) *fdCliFormat {
	return &fdCliFormat{
		stringFormatWriter: string_writer_format.MakeColor[string](
			co,
			relativePathStringFormatWriter,
			string_writer_format.ColorTypePointer,
		),
	}
}

func (f *fdCliFormat) WriteStringFormat(
	w io.StringWriter,
	k *FD,
) (n int64, err error) {
	// TODO-P2 add abbreviation

	var n1 int64

	n1, err = f.stringFormatWriter.WriteStringFormat(w, k.String())
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type kennungCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeKennungCliFormat(
	co string_writer_format.ColorOptions,
) *kennungCliFormat {
	return &kennungCliFormat{
		stringFormatWriter: string_writer_format.MakeColor[string](
			co,
			string_writer_format.MakeString[string](),
			string_writer_format.ColorTypePointer,
		),
	}
}

func (f *kennungCliFormat) WriteStringFormat(
	w io.StringWriter,
	k KennungPtr,
) (n int64, err error) {
	// TODO-P2 add abbreviation

	parts := k.Parts()

	var n1 int64

	n1, err = f.stringFormatWriter.WriteStringFormat(w, parts[0])
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var n2 int
	n2, err = w.WriteString(parts[1])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = f.stringFormatWriter.WriteStringFormat(w, parts[2])
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type typCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeTypCliFormat(co string_writer_format.ColorOptions) *typCliFormat {
	return &typCliFormat{
		stringFormatWriter: string_writer_format.MakeColor[string](
			co,
			string_writer_format.MakeString[string](),
			string_writer_format.ColorTypeType,
		),
	}
}

func (f *typCliFormat) WriteStringFormat(
	w io.StringWriter,
	k *Typ,
) (n int64, err error) {
	v := k.String()

	return f.stringFormatWriter.WriteStringFormat(w, v)
}

type etikettenCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeEtikettenCliFormat() *etikettenCliFormat {
	return &etikettenCliFormat{
		stringFormatWriter: string_writer_format.MakeString[string](),
	}
}

func (f *etikettenCliFormat) WriteStringFormat(
	w io.StringWriter,
	k EtikettSet,
) (n int64, err error) {
	v := iter.StringDelimiterSeparated[Etikett](k, " ")

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
