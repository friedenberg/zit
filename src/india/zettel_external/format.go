package zettel_external

import "github.com/friedenberg/zit/src/charlie/gattung"

type Parser = gattung.Parser[Zettel, *Zettel]
type Formatter = gattung.Formatter[Zettel, *Zettel]

type Format interface {
	Parser
	Formatter
}
