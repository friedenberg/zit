package zettel_external

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type Parser = schnittstellen.Parser[Zettel, *Zettel]
type Formatter = schnittstellen.Formatter[Zettel, *Zettel]

type Format interface {
	Parser
	Formatter
}
