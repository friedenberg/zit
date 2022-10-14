package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

type AkteWriterFactory interface {
	AkteWriter() (sha.WriteCloser, error)
}

type AkteReaderFactory interface {
	AkteReader(sha.Sha) (io.ReadCloser, error)
}

type FormatContextRead struct {
	Zettel            Zettel
	AktePath          string
	In                io.Reader
	RecoverableErrors errors.ErrorMulti
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
