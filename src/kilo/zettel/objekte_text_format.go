package zettel

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/india/erworben"
)

type objekteTextFormat struct {
	objekteTextParser
	objekteTextFormatter
}

func MakeObjekteTextFormat(
	akteFactory gattung.AkteIOFactory,
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
