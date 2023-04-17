package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/erworben"
)

type objekteTextFormat struct {
	metadatei.TextParser
	objekteTextFormatter
}

func MakeObjekteTextFormat(
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
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
