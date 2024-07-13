package object_metadata

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

type textFormatter struct {
	textFormatterCommon
	sequence []interfaces.FuncWriterElementInterface[TextFormatterContext]
}

func MakeTextFormatterMetadateiAktePath(
	options TextFormatterOptions,
	blobReaderFactory interfaces.BlobReaderFactory,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	if blobReaderFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		blobFactory:          blobReaderFactory,
		blobFormatter:        blobFormatter,
		TextFormatterOptions: options,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []interfaces.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadataFormat,
			common.writePathType,
			common.writeComments,
			common.writeBoundary,
		},
	}
}

func MakeTextFormatterMetadataOnly(
	options TextFormatterOptions,
	blobReaderFactory interfaces.BlobReaderFactory,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	if blobReaderFactory == nil {
		panic("akte reader factory is nil")
	}

	common := textFormatterCommon{
		blobFactory:          blobReaderFactory,
		blobFormatter:        blobFormatter,
		TextFormatterOptions: options,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []interfaces.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadataFormat,
			common.writeShaTyp,
			common.writeComments,
			common.writeBoundary,
		},
	}
}

func MakeTextFormatterMetadataInlineBlob(
	options TextFormatterOptions,
	blobReaderFactory interfaces.BlobReaderFactory,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	if blobReaderFactory == nil {
		panic("blob reader factory is nil")
	}

	common := textFormatterCommon{
		blobFactory:          blobReaderFactory,
		blobFormatter:        blobFormatter,
		TextFormatterOptions: options,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []interfaces.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBoundary,
			common.writeCommonMetadataFormat,
			common.writeTyp,
			common.writeComments,
			common.writeBoundary,
			common.writeNewLine,
			common.writeAkte,
		},
	}
}

func MakeTextFormatterExcludeMetadata(
	options TextFormatterOptions,
	blobFactory interfaces.BlobReaderFactory,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	if blobFactory == nil {
		panic("blob reader factory is nil")
	}

	common := textFormatterCommon{
		blobFactory:          blobFactory,
		blobFormatter:        blobFormatter,
		TextFormatterOptions: options,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []interfaces.FuncWriterElementInterface[TextFormatterContext]{
			common.writeAkte,
		},
	}
}

func (f textFormatter) FormatMetadata(
	w io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteSeq(w, c, f.sequence...)
}
