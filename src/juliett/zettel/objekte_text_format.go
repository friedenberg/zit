package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type objekteTextFormat struct {
	metadatei.TextFormatter
	metadatei.TextParser
}

func MakeObjekteTextFormat(
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter script_config.RemoteScript,
) *objekteTextFormat {
	return &objekteTextFormat{
		TextParser: metadatei.MakeTextParser(
			akteFactory,
			akteFormatter,
		),
		TextFormatter: metadatei.MakeTextFormatterMetadateiOnly(
			akteFactory,
			akteFormatter,
		),
	}
}
