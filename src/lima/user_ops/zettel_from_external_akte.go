package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/collections"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
	"github.com/friedenberg/zit/zettel_transacted"
)

type ZettelFromExternalAkte struct {
	Umwelt    *umwelt.Umwelt
	Etiketten etikett.Set
	Delete    bool
}

func (c ZettelFromExternalAkte) Run(args ...string) (results collections.SetTransacted, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	results = collections.MakeSetUniqueTransacted(len(args))

	for _, arg := range args {
		var z zettel.Zettel

		if z, err = c.zettelForAkte(store, arg); err != nil {
			err = errors.Error(err)
			return
		}

		var tz zettel_transacted.Transacted

		if tz, err = store.StoreObjekten().Create(z); err != nil {
			err = errors.Error(err)
			return
		}

		results.Add(tz)

		if c.Delete {
			if err = os.Remove(arg); err != nil {
				err = errors.Error(err)
				return
			}

			stdprinter.Errf("[%s] (deleted)\n", arg)
		}

		//TODO-P3,D3 only emit if created rather than refound
		stdprinter.Outf("%s (created)\n", tz.Named)
	}

	return
}

func (c ZettelFromExternalAkte) zettelForAkte(store store_with_lock.Store, aktePath string) (z zettel.Zettel, err error) {
	z.Etiketten = c.Etiketten

	var akteWriter age_io.Writer

	if akteWriter, err = store.StoreObjekten().AkteWriter(); err != nil {
		err = errors.Error(err)
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(aktePath); err != nil {
		err = errors.Error(err)
		return
	}

	defer open_file_guard.Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Error(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = z.Bezeichnung.Set(path.Base(aktePath)); err != nil {
		err = errors.Error(err)
		return
	}

	z.Akte = akteWriter.Sha()

	if err = z.Typ.Set(path.Ext(aktePath)); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
