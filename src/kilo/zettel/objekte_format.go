package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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

type ObjekteParser = schnittstellen.Parser[
	ObjekteParserContext,
	*ObjekteParserContext,
]

type ObjekteFormatter = schnittstellen.Formatter[
	ObjekteFormatterContext,
	*ObjekteFormatterContext,
]

type ObjekteFormat interface {
	ObjekteParser
	ObjekteFormatter
}
