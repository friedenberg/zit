package zettel

import (
	"io"

	"github.com/friedenberg/zit/charlie/sha"
	age_io "github.com/friedenberg/zit/delta/age_io"
)

type AkteWriterFactory interface {
	AkteWriter() (age_io.Writer, error)
}

type AkteReaderFactory interface {
	AkteReader(sha.Sha) (io.ReadCloser, error)
}

type FormatContextRead struct {
	Zettel           Zettel
	AktePath         string
	In               io.Reader
	RecoverableError error
	AkteWriterFactory
}

type FormatContextWrite struct {
	Zettel           Zettel
	Out              io.Writer
	IncludeAkte      bool
	ExternalAktePath string
	AkteReaderFactory
}

type Format interface {
	ReadFrom(*FormatContextRead) (int64, error)
	WriteTo(FormatContextWrite) (int64, error)
}
