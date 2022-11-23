package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

type AkteWriterFactory = metadatei_io.AkteWriterFactory
type AkteReaderFactory = metadatei_io.AkteReaderFactory

type FormatContextRead struct {
	Zettel            Zettel
	AktePath          string
	In                io.Reader
	RecoverableErrors errors.Multi
	AkteWriterFactory
}

type FormatContextWrite struct {
	Zettel           Zettel
	Out              io.Writer
	IncludeAkte      bool
	FormatScript     konfig.RemoteScript
	ExternalAktePath string
	AkteReaderFactory
}

type Format interface {
	ReadFrom(*FormatContextRead) (int64, error)
	WriteTo(FormatContextWrite) (int64, error)
}
