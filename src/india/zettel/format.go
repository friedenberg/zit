package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/india/konfig"
)

type FormatContextRead struct {
	Zettel            Objekte
	AktePath          string
	In                io.Reader
	RecoverableErrors errors.Multi
	gattung.AkteWriterFactory
}

type FormatContextWrite struct {
	Zettel           Objekte
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

type FormatToFormat2 struct {
	format.Format[Objekte]
}

func (ftf *FormatToFormat2) ReadFrom(
	c *FormatContextRead,
) (n int64, err error) {
	//TODO
	// AktePath          string
	// In                io.Reader
	// RecoverableErrors errors.Multi
	// gattung.AkteWriterFactory

	if n, err = ftf.Format.ReadFormat(c.In, &c.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ftf *FormatToFormat2) WriteTo(
	c FormatContextWrite,
) (n int64, err error) {
	//TODO
	// IncludeAkte      bool
	// FormatScript     konfig.RemoteScript
	// ExternalAktePath string
	// gattung.AkteReaderFactory

	if n, err = ftf.Format.WriteFormat(c.Out, &c.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
