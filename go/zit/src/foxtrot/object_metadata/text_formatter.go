package object_metadata

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
)

type textFormatter struct {
	textFormatterCommon
	sequence []interfaces.FuncWriterElementInterface[TextFormatterContext]
}

func MakeTextFormatterMetadataBlobPath(
	fs_home fs_home.Home,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	common := textFormatterCommon{
		fs_home:       fs_home,
		blobFormatter: blobFormatter,
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
	fs_home fs_home.Home,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	common := textFormatterCommon{
		fs_home:       fs_home,
		blobFormatter: blobFormatter,
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
	fs_home fs_home.Home,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	common := textFormatterCommon{
		fs_home:       fs_home,
		blobFormatter: blobFormatter,
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
			common.writeBlob,
		},
	}
}

func MakeTextFormatterExcludeMetadata(
	fs_home fs_home.Home,
	blobFormatter script_config.RemoteScript,
) textFormatter {
	common := textFormatterCommon{
		fs_home:       fs_home,
		blobFormatter: blobFormatter,
	}

	return textFormatter{
		textFormatterCommon: common,
		sequence: []interfaces.FuncWriterElementInterface[TextFormatterContext]{
			common.writeBlob,
		},
	}
}

func (f textFormatter) FormatMetadata(
	w io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteSeq(w, c, f.sequence...)
}
