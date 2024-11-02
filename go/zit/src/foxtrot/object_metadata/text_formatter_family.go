package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type TextFormatterFamily struct {
	BlobPath     TextFormatter
	InlineBlob   TextFormatter
	MetadataOnly TextFormatter
	BlobOnly     TextFormatter
}

func MakeTextFormatterFamily(
	dirLayout dir_layout.DirLayout,
	formatter script_config.RemoteScript,
) TextFormatterFamily {
	return TextFormatterFamily{
		BlobPath:     MakeTextFormatterMetadataBlobPath(dirLayout, formatter),
		InlineBlob:   MakeTextFormatterMetadataInlineBlob(dirLayout, formatter),
		MetadataOnly: MakeTextFormatterMetadataOnly(dirLayout, formatter),
		BlobOnly:     MakeTextFormatterExcludeMetadata(dirLayout, formatter),
	}
}
