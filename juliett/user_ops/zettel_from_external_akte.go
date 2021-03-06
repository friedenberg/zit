package user_ops

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/delta/umwelt"
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

	results.SetNamed = make(map[string]_NamedZettel, len(args))

	for _, arg := range args {
		var z _Zettel

		if z, err = c.zettelForAkte(store, arg); err != nil {
			err = _Error(err)
			return
		}

		var named _NamedZettel

		if named, err = store.Zettels().Create(z); err != nil {
			err = _Error(err)
			return
		}

		results.SetNamed[named.Hinweis.String()] = named

		if c.Delete {
			if err = os.Remove(arg); err != nil {
				err = _Error(err)
				return
			}

			stdprinter.Errf("[%s] (deleted)\n", arg)
		}

		//TODO-P3,D3 only emit if created rather than refound
		stdprinter.Outf("[%s %s] (created)\n", named.Hinweis, named.Sha)
	}

	return
}

func (c ZettelFromExternalAkte) zettelForAkte(store store_with_lock.Store, aktePath string) (z _Zettel, err error) {
	z.Etiketten = c.Etiketten

	var akteWriter objekte.Writer

	if akteWriter, err = store.Zettels().AkteWriter(); err != nil {
		err = _Error(err)
		return
	}

	var f *os.File

	if f, err = open_file_guard.Open(aktePath); err != nil {
		err = _Error(err)
		return
	}

	defer open_file_guard.Close(f)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = _Error(err)
		return
	}

	if err = akteWriter.Close(); err != nil {
		err = _Error(err)
		return
	}

	if err = z.Bezeichnung.Set(path.Base(aktePath)); err != nil {
		err = _Error(err)
		return
	}

	z.Akte = akteWriter.Sha()

	if err = z.AkteExt.Set(path.Ext(aktePath)); err != nil {
		err = _Error(err)
		return
	}

	return
}
