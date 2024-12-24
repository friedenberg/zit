package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type TextFormat struct {
	TextFormatterFamily
	TextParser
}

func MakeTextFormat(
	dirLayout dir_layout.DirLayout,
	blobFormatter script_config.RemoteScript,
) TextFormat {
	return TextFormat{
		TextParser: MakeTextParser(
			dirLayout,
			blobFormatter,
		),
		TextFormatterFamily: MakeTextFormatterFamily(
			dirLayout,
			blobFormatter,
		),
	}
}
