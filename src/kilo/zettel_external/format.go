package zettel_external

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type (
	Parser    = objekte.AkteParser[*Zettel]
	Formatter = schnittstellen.Formatter[Zettel, *Zettel]
)

type Format interface {
	Parser
	Formatter
}
