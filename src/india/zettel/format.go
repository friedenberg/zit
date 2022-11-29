package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/echo/konfig"
)

type FormatContextRead struct {
	Zettel            Zettel
	AktePath          string
	In                io.Reader
	RecoverableErrors errors.Multi
	gattung.AkteWriterFactory
}

type FormatContextWrite struct {
	Zettel           Zettel
	Out              io.Writer
	IncludeAkte      bool
	FormatScript     konfig.RemoteScript
	ExternalAktePath string
	gattung.AkteReaderFactory
}

type Format interface {
	ReadFrom(*FormatContextRead) (int64, error)
	WriteTo(FormatContextWrite) (int64, error)
}
