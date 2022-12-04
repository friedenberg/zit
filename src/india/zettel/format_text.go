package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

type Text struct {
	DoNotWriteEmptyBezeichnung bool
	TypError                   error
}

func (f Text) ReadFrom(c *FormatContextRead) (n int64, err error) {
	state := &FormatMetadateiText{
		context: c,
	}

	var akteWriter sha.WriteCloser

	if c.AkteWriterFactory == nil {
		err = errors.Errorf("akte writer factory is nil")
		return
	}

	if akteWriter, err = c.AkteWriter(); err != nil {
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
