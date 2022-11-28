package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/hotel/cwd_files"
)

func (s *Store) CheckinTyp(p string) (t *typ.Transacted, err error) {
	return
}

func (s *Store) WriteTyp(t *typ.Transacted) (te *typ.External, err error) {
	te = &typ.External{
		FD: cwd_files.File{
			Path: fmt.Sprintf("%s.%s", t.Kennung(), s.Konfig.Transacted.Objekte.Akte.TypFileExtension),
		},
		//TODO move to central place
		Objekte: t.Objekte,
		Sha:     t.Objekte.Sha,
		Kennung: t.Sku.Kennung,
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(te.FD.Path); err != nil {
		if errors.IsExist(err) {
			err = s.ReadTyp(te)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, f.Close)

	format := typ.MakeFormatText(s.storeObjekten)

	if _, err = format.WriteFormat(f, &te.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTyp(t *typ.External) (err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	var f *os.File

	if f, err = files.Open(t.FD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	if _, err = format.ReadFormat(f, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
