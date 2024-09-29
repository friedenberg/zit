package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
)

type TextFormatterFamily struct {
	BlobPath     TextFormatter
	InlineBlob   TextFormatter
	MetadataOnly TextFormatter
	BlobOnly     TextFormatter
}

func MakeTextFormatterFamily(
	fs_home fs_home.Home,
	formatter script_config.RemoteScript,
) TextFormatterFamily {
	return TextFormatterFamily{
		BlobPath:     MakeTextFormatterMetadataBlobPath(fs_home, formatter),
		InlineBlob:   MakeTextFormatterMetadataInlineBlob(fs_home, formatter),
		MetadataOnly: MakeTextFormatterMetadataOnly(fs_home, formatter),
		BlobOnly:     MakeTextFormatterExcludeMetadata(fs_home, formatter),
	}
}
