package zettel

import (
	"bufio"
	"io"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type textStateReadField int

const (
	textStateReadFieldEmpty = textStateReadField(iota)
	textStateReadFieldFirstBoundary
	textStateReadFieldSecondBoundary
	textStateReadFieldAkteBody
)

type textStateReadAkte int

const (
	// no akte file or ext, therefore it's inline
	// yes akte just ext and it's an inline type, therefore it's inline
	// yes akte just ext and it's not an inline type, therefore error
	// yes akte file and ext, therefore it's a file
	// yes akte file and ext and content inline, therefore error
	textStateReadAkteInline            = textStateReadAkte(iota)
	textStateReadAkteFileWithExtension = textStateReadAkte(iota)
	textStateReadAkteJustExtension     = textStateReadAkte(iota)
)

type textStateRead struct {
	textStateReadMetadatei
	context          *FormatContextRead
	field            textStateReadField
	didReadAkte      bool
	metadataiAkteSha sha.Sha
	readAkteSha      sha.Sha
	akteWriter       sha.WriteCloser
}

func (s *textStateRead) Close() (err error) {
	if s.akteWriter != nil {
		if err = s.akteWriter.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.didReadAkte = true
		s.readAkteSha = s.akteWriter.Sha()
	}

	s.context.RecoverableErrors = s.recoverableErrors

	if err = s.textStateReadMetadatei.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.context.Zettel = *s.textStateReadMetadatei.Zettel

	return
}

type textStateWrite struct {
	Zettel
}

type textStateReadMetadatei struct {
	*Zettel
	etiketten         kennung.EtikettMutableSet
	aktePath          string
	recoverableErrors errors.Multi
}

func (s textStateReadMetadatei) Close() (err error) {
	s.Zettel.Etiketten = s.etiketten.Copy()
	return
}

func (f textStateReadMetadatei) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]format.FuncReadLine{
					"#": f.Zettel.Bezeichnung.Set,
					"%": format.MakeLineReaderNop(),
					"-": f.etiketten.AddString,
					"!": f.readTyp,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f textStateReadMetadatei) readTyp(desc string) (err error) {
	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	// path
	if files.Exists(desc) {
		errors.Log().Print("valid path", desc)

		if err = f.Zettel.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		f.aktePath = desc

		return
	}

	//TODO handl akte descs that are invalid files

	shaError := f.Akte.Set(head)

	if tail == "" {
		//sha or ext
		if shaError != nil {
			if err = f.Zettel.Typ.Set(head); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	} else {
		//sha.ext or error
		if shaError != nil {
			f.recoverableErrors.Add(
				errors.Wrap(
					ErrHasInvalidAkteShaOrFilePath{
						Value: head,
					},
				),
			)

			return
		}

		if err = f.Zettel.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
