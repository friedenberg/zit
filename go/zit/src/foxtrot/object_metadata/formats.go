package object_metadata

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

type TextFormat struct {
	TextFormatter
	TextParser
}

func MakeTextFormat(
	blobIOFactory interfaces.BlobIOFactory,
	blobFormatter script_config.RemoteScript,
) TextFormat {
	return TextFormat{
		TextParser: MakeTextParser(
			blobIOFactory,
			blobFormatter,
		),
		TextFormatter: MakeTextFormatterMetadateiOnly(
			TextFormatterOptions{},
			blobIOFactory,
			blobFormatter,
		),
	}
}
