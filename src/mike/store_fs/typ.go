package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/typ_toml"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/hotel/cwd_files"
)

func (s *Store) CheckinTyp(p string) (t *typ.Transacted, err error) {
	return
}

func (s *Store) WriteTyp(t *typ.Transacted) (te *typ.External, err error) {
	te = &typ.External{
		FD: cwd_files.File{
			Path: fmt.Sprintf("%s.%s", t.Kennung(), s.konfig.FileExtensions.Typ),
		},
		//TODO-P2 move to central place
		Objekte: t.Objekte,
		Sku: sku.External[kennung.Typ, *kennung.Typ]{
			Sha:     t.ObjekteSha(),
			Kennung: t.Sku.Kennung,
		},
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

	format := typ_toml.MakeFormatText(s.storeObjekten)

	if _, err = format.WriteFormat(f, &te.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTyp(t *typ.External) (err error) {
	format := typ_toml.MakeFormatText(s.storeObjekten)

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
