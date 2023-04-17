package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
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

type ObjekteParser = schnittstellen.ParserInterface[metadatei.ParserContext]

type ObjekteFormatter = schnittstellen.Formatter[
	ObjekteFormatterContext,
	*ObjekteFormatterContext,
]

type ObjekteFormat interface {
	ObjekteParser
	ObjekteFormatter
}

func (c *ObjekteParserContext) GetMetadateiPtr() *metadatei.Metadatei {
	return &c.Zettel.Metadatei
}

func (c *ObjekteParserContext) SetAkteFD(fd kennung.FD) (err error) {
	// TODO read into akte writer?
	c.AktePath = fd.Path
	c.Zettel.Metadatei.AkteSha = fd.Sha
	return
}
