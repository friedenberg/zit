package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type objekteTextFormat struct {
	metadatei.TextParser
	objekteTextFormatter
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
		objekteTextFormatter: objekteTextFormatter{
			AkteFactory: akteFactory,
		},
	}
}
