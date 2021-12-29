package zettel

import "io"

type AkteWriterFactory interface {
	AkteWriter() (_ObjekteWriter, error)
}

type AkteReaderFactory interface {
	AkteReader(_Sha) (io.ReadCloser, error)
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
