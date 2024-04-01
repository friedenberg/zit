package metadatei

import (
	"io"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/src/charlie/script_config"
)

type textFormatter struct {
	textFormatterCommon
	sequence []schnittstellen.FuncWriterElementInterface[TextFormatterContext]
}

func MakeTextFormatterMetadateiAktePath(
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
			common.writePathTyp,
			common.writeComments,
			common.writeBoundary,
		},
	}
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
			common.writeComments,
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
			common.writeComments,
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

func (f textFormatter) FormatMetadatei(
	w io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteSeq(w, c, f.sequence...)
}
