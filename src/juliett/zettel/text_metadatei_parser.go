package zettel

import (
	"bufio"
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type textMetadateiParser struct {
	awf schnittstellen.AkteWriterFactory
	// context *FormatContextRead

	aktePath string
	akteSha  sha.Sha
}

func MakeTextMetadateiParser(
	awf schnittstellen.AkteWriterFactory,
) textMetadateiParser {
	return textMetadateiParser{
		awf: awf,
	}
}

func (f *textMetadateiParser) Parse(r1 io.Reader, z *Objekte) (n int64, err error) {
	return f.ReadFormat(r1, z)
}

func (f *textMetadateiParser) ReadFormat(r1 io.Reader, z *Objekte) (n int64, err error) {
	etiketten := kennung.MakeEtikettMutableSet()

	defer func() {
		z.Etiketten = etiketten.Copy()
	}()

	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					"#": z.Bezeichnung.Set,
					"%": format.MakeLineReaderNop(),
					"-": etiketten.StringAdder(),
					"!": func(v string) (err error) {
						return f.readTyp(z, v)
					},
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *textMetadateiParser) readTyp(
	z *Objekte,
	desc string,
) (err error) {
	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	errors.TodoP3("handle akte descs that are invalid files")
	//! <path>.<typ ext>
	switch {
	case files.Exists(desc):
		if err = z.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		f.aktePath = desc

		var akteWriter sha.WriteCloser

		if akteWriter, err = f.awf.AkteWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if akteWriter == nil {
			err = errors.Errorf("akte writer is nil")
			return
		}

		var fi *os.File

		if fi, err = os.Open(desc); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, fi.Close)

		if _, err = io.Copy(akteWriter, fi); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = akteWriter.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}

		f.akteSha = sha.Make(akteWriter.Sha())

	//! <sha>.<typ ext>
	case tail != "":
		if err = z.Akte.Set(head); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = z.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>
	case tail == "":
		if err = z.Akte.Set(head); err == nil {
			return
		}

		err = nil

		fallthrough

	//! <typ ext>
	default:
		if err = z.Typ.Set(head); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
