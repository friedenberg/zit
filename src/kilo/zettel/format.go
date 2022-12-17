package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/konfig"
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
	FormatScript     konfig.RemoteScript
	ExternalAktePath string
}

type Format interface {
	ReadFrom(*FormatContextRead) (int64, error)
	WriteTo(FormatContextWrite) (int64, error)
}
