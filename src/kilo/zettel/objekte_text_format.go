package zettel

import (
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/schnittstellen"
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
