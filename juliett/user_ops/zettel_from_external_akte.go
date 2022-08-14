package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/age_io"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type ZettelFromExternalAkte struct {
	Umwelt    *umwelt.Umwelt
	Etiketten etikett.Set
	Delete    bool
}

func (c ZettelFromExternalAkte) Run(args ...string) (results ZettelResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	results.SetNamed = make(map[string]stored_zettel.Named, len(args))

	for _, arg := range args {
		var z zettel.Zettel

		if z, err = c.zettelForAkte(store, arg); err != nil {
			err = errors.Error(err)
			return
		}

		var tz stored_zettel.Transacted

		if tz, err = store.Zettels().Create(z); err != nil {
			err = errors.Error(err)
			return
		}

		results.SetNamed[tz.Hinweis.String()] = tz.Named

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

	var akteWriter objekte.Writer

	if akteWriter, err = store.Zettels().AkteWriter(); err != nil {
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

	if err = z.AkteExt.Set(path.Ext(aktePath)); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
