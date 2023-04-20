package metadatei

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/ohio"
	"github.com/friedenberg/zit/src/bravo/script_config"
)

type textFormatter struct {
	textFormatterCommon
	sequence []schnittstellen.FuncWriterElementInterface[TextFormatterContext]
}

func MakeTextFormatterMetadateiOnly(
	akteFactory schnittstellen.AkteReaderFactory,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	if akteFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		akteFactory:   akteFactory,
		akteFormatter: akteFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []schnittstellen.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadateiFormat,
			common.writeShaTyp,
			common.writeBoundary,
		},
	}
}

func MakeTextFormatterMetadateiInlineAkte(
	akteFactory schnittstellen.AkteReaderFactory,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	if akteFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		akteFactory:   akteFactory,
		akteFormatter: akteFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []schnittstellen.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadateiFormat,
			common.writeTyp,
			common.writeBoundary,
			common.writeNewLine,
			common.writeAkte,
		},
	}
}

func MakeTextFormatterExcludeMetadatei(
	akteFactory schnittstellen.AkteReaderFactory,
	akteFormatter script_config.RemoteScript,
) textFormatter {
	if akteFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		akteFactory:   akteFactory,
		akteFormatter: akteFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []schnittstellen.FuncWriterElementInterface[TextFormatterContext]{
			common.writeAkte,
		},
	}
}

func (f textFormatter) Format(
	w io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteSeq(w, c, f.sequence...)
}
