package metadatei

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
)

type TextFormat struct {
	TextFormatter
	TextParser
}

func MakeTextFormat(
	akteFactory schnittstellen.AkteIOFactory,
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
