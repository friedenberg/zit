package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/erworben"
)

type objekteTextFormat struct {
	objekteTextParser
	objekteTextFormatter
}

func MakeObjekteTextFormat(
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) *objekteTextFormat {
	return &objekteTextFormat{
		objekteTextParser: MakeObjekteTextParser(
			akteFactory,
			akteFormatter,
		),
		objekteTextFormatter: objekteTextFormatter{
			AkteFactory: akteFactory,
		},
	}
}
