package metadatei

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

type TextFormat struct {
	TextFormatter
	TextParser
}

func MakeTextFormat(
	akteFactory interfaces.BlobIOFactory,
	akteFormatter script_config.RemoteScript,
) TextFormat {
	return TextFormat{
		TextParser: MakeTextParser(
			akteFactory,
			akteFormatter,
		),
		TextFormatter: MakeTextFormatterMetadateiOnly(
			TextFormatterOptions{},
			akteFactory,
			akteFormatter,
		),
	}
}
