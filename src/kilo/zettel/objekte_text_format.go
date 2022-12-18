package zettel

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/india/konfig"
)

type objekteTextFormat struct {
	textParser
}

func MakeObjekteTextFormat(
	akteFactory gattung.AkteIOFactory,
	akteFormatter konfig.RemoteScript,
) *objekteTextFormat {
	return &objekteTextFormat{
		textParser: MakeTextParser(
			akteFactory,
			akteFormatter,
		),
	}
}
