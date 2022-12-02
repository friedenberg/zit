package zettel

import (
	"bufio"
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type FormatMetadateiText struct {
	recoverableErrors errors.Multi

	context *FormatContextRead

	etiketten kennung.EtikettMutableSet

	aktePath string
	akteSha  sha.Sha
}

func (s *FormatMetadateiText) close() (err error) {
	s.context.RecoverableErrors = s.recoverableErrors
	s.context.Zettel.Etiketten = s.etiketten.Copy()

	return
}

func (f FormatMetadateiText) ReadFrom(r1 io.Reader) (n int64, err error) {
	defer errors.Deferred(&err, f.close)

	f.etiketten = kennung.MakeEtikettMutableSet()

	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]format.FuncReadLine{
					"#": f.context.Zettel.Bezeichnung.Set,
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

func (f FormatMetadateiText) readTyp(desc string) (err error) {
	if desc == "" {
		return
	}

	tail := path.Ext(desc)
	head := desc[:len(desc)-len(tail)]

	//TODO handl akte descs that are invalid files
	//! <path>.<typ ext>
	switch {
	case files.Exists(desc):
		if err = f.context.Zettel.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

		f.aktePath = desc

		var akteWriter sha.WriteCloser

		if akteWriter, err = f.context.AkteWriter(); err != nil {
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
		if err = f.context.Zettel.Akte.Set(head); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = f.context.Zettel.Typ.Set(tail); err != nil {
			err = errors.Wrap(err)
			return
		}

	//! <sha>
	case tail == "":
		if err = f.context.Zettel.Akte.Set(head); err == nil {
			return
		}

		err = nil

		fallthrough

	//! <typ ext>
	default:
		if err = f.context.Zettel.Typ.Set(head); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
