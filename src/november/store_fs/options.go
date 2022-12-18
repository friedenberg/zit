package store_fs

import (
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

type OptionsReadExternal struct {
	Parser   zettel.ObjekteParser //TODO-P1 switch to zettel_external.Parser
	Zettelen map[hinweis.Hinweis]zettel.Transacted
}

type CheckoutOptions struct {
	Force bool
	CheckoutMode
	Formatter zettel.ObjekteFormatter //TODO-P1 switch to zettel_external.Formatter
	Zettelen  map[hinweis.Hinweis]zettel.Transacted
}
