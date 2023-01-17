package zettel

import (
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type ObjekteParserContext struct {
	Zettel   Objekte
	AktePath string
	Errors   error
}

type ObjekteFormatterContext struct {
	Zettel           Objekte
	IncludeAkte      bool
	ExternalAktePath string
}

type ObjekteParser = gattung.Parser[
	ObjekteParserContext,
	*ObjekteParserContext,
]

type ObjekteFormatter = gattung.Formatter[
	ObjekteFormatterContext,
	*ObjekteFormatterContext,
]

type ObjekteFormat interface {
	ObjekteParser
	ObjekteFormatter
}
