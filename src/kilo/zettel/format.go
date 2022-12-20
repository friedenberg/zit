package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type FormatContextRead struct {
	Zettel            Objekte
	AktePath          string
	RecoverableErrors errors.Multi
}

type FormatContextWrite struct {
	Zettel           Objekte
	IncludeAkte      bool
	ExternalAktePath string
}

type Format interface {
	Parse(io.Reader, *FormatContextRead) (int64, error)
	Format(io.Writer, FormatContextWrite) (int64, error)
}

type ObjekteParser = gattung.Parser[FormatContextRead, *FormatContextRead]
type ObjekteFormatter = gattung.Formatter[FormatContextWrite, *FormatContextWrite]

type ObjekteFormat interface {
	ObjekteParser
	ObjekteFormatter
}
