package zettel

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/india/konfig"
)

type objekteTextFormat struct {
	objekteTextParser
}

func MakeObjekteTextFormat(
	akteFactory gattung.AkteIOFactory,
	akteFormatter konfig.RemoteScript,
) *objekteTextFormat {
	return &objekteTextFormat{
		objekteTextParser: MakeTextParser(
			akteFactory,
			akteFormatter,
		),
	}
}
