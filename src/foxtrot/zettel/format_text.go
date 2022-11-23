package zettel

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
)

const (
	MetadateiBoundary = "---"
)

type Text struct {
	DoNotWriteEmptyBezeichnung bool
	TypError                   error
}

func (f Text) ReadFrom(c *FormatContextRead) (n int64, err error) {
	state := &textStateRead{
		textStateReadMetadatei: textStateReadMetadatei{
			Zettel:    &Zettel{},
			etiketten: etikett.MakeMutableSet(),
		},
		context: c,
	}

	defer func() {
		err1 := state.Close()

		if err == nil {
			err = err1
		}

		if !state.context.Zettel.Akte.IsNull() {
			return
		}

		//TODO log the following states
		if !state.metadataiAkteSha.IsNull() {
			state.context.Zettel.Akte = state.metadataiAkteSha
			return
		}

		if !state.readAkteSha.IsNull() {
			state.context.Zettel.Akte = state.readAkteSha
			return
		}
	}()

	if c.AkteWriterFactory == nil {
		err = errors.Errorf("akte writer factory is nil")
		return
	}

	if state.akteWriter, err = c.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if state.akteWriter == nil {
		err = errors.Errorf("akte writer is nil")
		return
	}

	mr := metadatei_io.Reader{
		Metadatei: state.textStateReadMetadatei,
		Akte:      state.akteWriter,
	}

	if n, err = mr.ReadFrom(c.In); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO move this to parallelized metadatei
	if state.aktePath != "" {
		var f *os.File

		if f, err = files.Open(state.aktePath); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer files.Close(f)

		if _, err = io.Copy(state.akteWriter, f); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
