package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type ObjekteParserContext struct {
	Zettel            Objekte
	AktePath          string
	RecoverableErrors errors.Multi
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
