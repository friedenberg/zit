package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/metadatei_io"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/india/erworben"
)

type objekteTextParser struct {
	AkteFactory                schnittstellen.AkteIOFactory
	AkteFormatter              erworben.RemoteScript
	DoNotWriteEmptyBezeichnung bool
	TypError                   error
}

func MakeObjekteTextParser(
	akteFactory schnittstellen.AkteIOFactory,
	akteFormatter erworben.RemoteScript,
) objekteTextParser {
	return objekteTextParser{
		AkteFactory:   akteFactory,
		AkteFormatter: akteFormatter,
	}
}

func (f objekteTextParser) Parse(
	r io.Reader,
	c *ObjekteParserContext) (n int64, err error) {
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

	if n, err = mr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	inlineAkteSha := akteWriter.Sha()
	c.AktePath = state.aktePath

	switch {
	case state.akteSha.IsNull() && !inlineAkteSha.IsNull():
		c.Zettel.Akte = sha.Make(inlineAkteSha)

	case !state.akteSha.IsNull() && inlineAkteSha.IsNull():
		c.Zettel.Akte = state.akteSha

	case !state.akteSha.IsNull() && !inlineAkteSha.IsNull():
		err = ErrHasInlineAkteAndFilePath{
			External: externalFile{
				Sha:  state.akteSha,
				Path: state.aktePath,
			},
			InlineSha: sha.Make(inlineAkteSha),
			Objekte:   c.Zettel,
		}

		return
	}

	return
}
