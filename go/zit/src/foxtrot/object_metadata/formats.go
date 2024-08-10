package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
)

type TextFormat struct {
	TextFormatterFamily
	TextParser
}

func MakeTextFormat(
	fs_home fs_home.Home,
	blobFormatter script_config.RemoteScript,
) TextFormat {
	return TextFormat{
		TextParser: MakeTextParser(
			fs_home,
			blobFormatter,
		),
		TextFormatterFamily: MakeTextFormatterFamily(
			fs_home,
			blobFormatter,
		),
	}
}
