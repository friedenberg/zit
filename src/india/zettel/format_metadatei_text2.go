package zettel

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/sha"
)

type FormatMetadateiText2 struct {
	af gattung.AkteWriterFactory

	DoNotWriteEmptyBezeichnung bool

	aktePath string
	akteSha  sha.Sha
}

func (f FormatMetadateiText2) ReadFormat(r1 io.Reader, z *Objekte) (n int64, err error) {
	etiketten := kennung.MakeEtikettMutableSet()

	defer func() {
		z.Etiketten = etiketten.Copy()
	}()

	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]format.FuncReadLine{
					"#": z.Bezeichnung.Set,
					"%": format.MakeLineReaderNop(),
					"-": etiketten.AddString,
					"!": f.makeReadTypFunc(z),
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f FormatMetadateiText2) makeReadTypFunc(
	z *Objekte,
) format.FuncReadLine {
	return func(desc string) (err error) {
		if desc == "" {
			return
		}

		tail := path.Ext(desc)
		head := desc[:len(desc)-len(tail)]

		//TODO handl akte descs that are invalid files
		//! <path>.<typ ext>
		switch {
		case files.Exists(desc):
			if err = z.Typ.Set(tail); err != nil {
				err = errors.Wrap(err)
				return
			}

			f.aktePath = desc

			var akteWriter sha.WriteCloser

			if akteWriter, err = f.af.AkteWriter(); err != nil {
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

			f.akteSha = akteWriter.Sha()

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
}

func (f *FormatMetadateiText2) WriteFormat(
	w1 io.Writer,
	z *Objekte,
) (n int64, err error) {
	w := format.NewWriter()

	if z.Bezeichnung.String() != "" || !f.DoNotWriteEmptyBezeichnung {
		w.WriteLines(
			fmt.Sprintf("# %s", z.Bezeichnung),
		)
	}

	for _, e := range z.Etiketten.Sorted() {
		//TODO
		if e.IsEmpty() {
			continue
		}

		w.WriteFormat("- %s", e)
	}

	switch {
	//TODO log this state
	case z.Akte.IsNull() && z.Typ.String() == "":
		break

	case z.Akte.IsNull():
		w.WriteLines(
			fmt.Sprintf("! %s", z.Typ),
		)

	case z.Typ.String() == "":
		w.WriteLines(
			fmt.Sprintf("! %s", z.Akte),
		)

	default:
		w.WriteLines(
			fmt.Sprintf("! %s.%s", z.Akte, z.Typ),
		)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
