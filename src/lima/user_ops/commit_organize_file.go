package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/india/changes"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
	"github.com/friedenberg/zit/src/lima/zettel_named"
	"github.com/friedenberg/zit/zettel_transacted"
)

type CommitOrganizeFile struct {
	Umwelt *umwelt.Umwelt
}

type CommitOrganizeFileResults struct {
}

func (c CommitOrganizeFile) Run(a, b organize_text.Text) (results CommitOrganizeFileResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var cs changes.Changes

	if cs, err = changes.ChangesFrom(a, b); err != nil {
		err = errors.Error(err)
		return
	}

	logz.Printf("%#v", cs)

	if len(cs.Added) == 0 && len(cs.Removed) == 0 && len(cs.New) == 0 {
		stdprinter.Err("no changes")
		return
	}

	toUpdate := make(map[string]zettel_named.Zettel)

	addOrGetToZettelToUpdate := func(hString string) (z zettel_named.Zettel, err error) {
		var h hinweis.Hinweis

		if h, err = hinweis.Make(hString); err != nil {
			err = errors.Error(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			var tz zettel_transacted.Transacted

			if tz, err = store.StoreObjekten().Read(h); err != nil {
				err = errors.Error(err)
				return
			}

			z = tz.Named
		}

		return
	}

	addEtikettToZettel := func(hString string, e etikett.Etikett) (err error) {
		var z zettel_named.Zettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Error(err)
			return
		}

		z.Stored.Zettel.Etiketten.Add(e)
		toUpdate[z.Hinweis.String()] = z

		stdprinter.Errf("Added etikett '%s' to zettel '%s'\n", e, z.Hinweis)

		return
	}

	removeEtikettFromZettel := func(hString string, e etikett.Etikett) (err error) {
		var z zettel_named.Zettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Error(err)
			return
		}

		z.Stored.Zettel.Etiketten.RemovePrefixes(e)
		toUpdate[z.Hinweis.String()] = z

		stdprinter.Errf("Removed etikett '%s' from zettel '%s'\n", e, z.Hinweis)

		return
	}

	for _, c := range cs.Removed {
		var e etikett.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = errors.Error(err)
			return
		}

		if err = removeEtikettFromZettel(c.Key, e); err != nil {
			err = errors.Error(err)
			return
		}
	}
	for _, c := range cs.Added {
		var e etikett.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = errors.Error(err)
			return
		}

		if err = addEtikettToZettel(c.Key, e); err != nil {
			err = errors.Error(err)
			return
		}
	}

	for bez, etts := range cs.New {
		z := zettel.Zettel{
			Etiketten: etts,
		}

		if err = z.Bezeichnung.Set(bez); err != nil {
			err = errors.Error(err)
			return
		}

		if err = z.Typ.Set("md"); err != nil {
			err = errors.Error(err)
			return
		}

		if c.Umwelt.Konfig.DryRun {
			stdprinter.Outf("[%s] (would create)\n", z.Bezeichnung)
			continue
		}

		var tz zettel_transacted.Transacted

		if tz, err = store.StoreObjekten().Create(z); err != nil {
			err = errors.Errorf("failed to create zettel: %s", err)
			return
		}

		stdprinter.Outf("%s (created)\n", tz.Named)
	}

	for _, z := range toUpdate {
		if c.Umwelt.Konfig.DryRun {
			stdprinter.Outf("[%s] (would update)\n", z.Hinweis)
			continue
		}

		if _, err = store.StoreObjekten().Update(z.Hinweis, z.Stored.Zettel); err != nil {
			stdprinter.Errf("failed to update zettel: %s", err)
		}
	}

	return
}
