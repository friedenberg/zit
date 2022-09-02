package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/open_file_guard"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
)

type ZettelFromExternalAkte struct {
	Umwelt    *umwelt.Umwelt
	Etiketten etikett.Set
	Delete    bool
}

func (c ZettelFromExternalAkte) Run(args ...string) (results zettel_transacted.Set, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	results = zettel_transacted.MakeSetUnique(len(args))

	for _, arg := range args {
		var z zettel.Zettel

		if z, err = c.zettelForAkte(store, arg); err != nil {
			err = errors.Wrap(err)
			return
		}

		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().Create(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		results.Add(tz)

		if c.Delete {
			if err = os.Remove(arg); err != nil {
				err = errors.Wrap(err)
				return
			}

			errors.PrintErrf("[%s] (deleted)", arg)
		}

		//TODO-P3,D3 only emit if created rather than refound
		errors.PrintOutf("%s (created)", tz.Named)
	}

	return
}

func (c ZettelFromExternalAkte) zettelForAkte(store store_with_lock.Store, aktePath string) (z zettel.Zettel, err error) {
	z.Etiketten = c.Etiketten

	var akteWriter sha.WriteCloser

	if akteWriter, err = store.StoreObjekten().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(aktePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = z.Bezeichnung.Set(path.Base(aktePath)); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.Akte = akteWriter.Sha()

	if err = z.Typ.Set(path.Ext(aktePath)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
