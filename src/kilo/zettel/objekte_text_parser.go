package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/metadatei_io"
	"github.com/friedenberg/zit/src/india/konfig"
)

// TODO-P0 migrate to format.Format
type objekteTextParser struct {
	AkteFactory                gattung.AkteIOFactory
	AkteFormatter              konfig.RemoteScript
	DoNotWriteEmptyBezeichnung bool
	TypError                   error
}

// TODO-P1 rename dto MakeObjekteTextParser
func MakeTextParser(
	akteFactory gattung.AkteIOFactory,
	akteFormatter konfig.RemoteScript,
) objekteTextParser {
	return objekteTextParser{
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
	}
}

func (f objekteTextParser) Parse(r io.Reader, o *Objekte) (n int64, err error) {
	return f.ReadFormat(r, o)
}

func (f objekteTextParser) ReadFormat(r io.Reader, o *Objekte) (n int64, err error) {
	c := &FormatContextRead{
		In:     r,
		Zettel: *o,
	}

	if n, err = f.ReadFrom(c); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f objekteTextParser) ReadFrom(c *FormatContextRead) (n int64, err error) {
	state := MakeTextMetadateiParser(
		f.AkteFactory,
	)

	var akteWriter sha.WriteCloser

	if f.AkteFactory == nil {
		err = errors.Errorf("akte factory is nil")
		return
	}

	if akteWriter, err = f.AkteFactory.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if akteWriter == nil {
		err = errors.Errorf("akte writer is nil")
		return
	}

	mr := metadatei_io.Reader{
		Metadatei: format.MakeReaderFrom[Objekte](state.ReadFormat, &c.Zettel),
		Akte:      akteWriter,
	}

	if n, err = mr.ReadFrom(c.In); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	inlineAkteSha := akteWriter.Sha()
	c.AktePath = state.aktePath

	//TODO-P1 handle akte path
	switch {
	case state.akteSha.IsNull() && !inlineAkteSha.IsNull():
		c.Zettel.Akte = inlineAkteSha

	case !state.akteSha.IsNull() && inlineAkteSha.IsNull():
		c.Zettel.Akte = state.akteSha

	case !state.akteSha.IsNull() && !inlineAkteSha.IsNull():
		err = ErrHasInlineAkteAndFilePath{
			External: externalFile{
				Sha:  state.akteSha,
				Path: state.aktePath,
			},
			InlineSha: inlineAkteSha,
			Objekte:   c.Zettel,
		}

		return
	}

	return
}
