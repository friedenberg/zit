package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type FormatContextRead struct {
	Zettel            Objekte
	AktePath          string
	In                io.Reader
	RecoverableErrors errors.Multi
}

type FormatContextWrite struct {
	Zettel           Objekte
	Out              io.Writer
	IncludeAkte      bool
	ExternalAktePath string
}

type Format interface {
	ReadFrom(*FormatContextRead) (int64, error)
	WriteTo(FormatContextWrite) (int64, error)
}

type ObjekteParser = gattung.Parser[Objekte, *Objekte]
type ObjekteFormatter = gattung.Formatter[Objekte, *Objekte]

type ObjekteFormat interface {
	ObjekteParser
	ObjekteFormatter
}
