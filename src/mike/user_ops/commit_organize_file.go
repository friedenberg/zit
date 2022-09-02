package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/kilo/changes"
	"github.com/friedenberg/zit/src/lima/store_with_lock"
)

type CommitOrganizeFile struct {
	Umwelt *umwelt.Umwelt
}

type CommitOrganizeFileResults struct {
}

func (c CommitOrganizeFile) Run(a, b organize_text.Text) (results CommitOrganizeFileResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	var cs changes.Changes

	if cs, err = changes.ChangesFrom(a, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Printf("%#v", cs)

	if len(cs.Added) == 0 && len(cs.Removed) == 0 && len(cs.New) == 0 {
		errors.PrintErr("no changes")
		return
	}

	toUpdate := make(map[string]zettel_named.Zettel)

	addOrGetToZettelToUpdate := func(hString string) (z zettel_named.Zettel, err error) {
		var h hinweis.Hinweis

		if h, err = hinweis.Make(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			var tz zettel_transacted.Zettel

			if tz, err = store.StoreObjekten().Read(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			z = tz.Named
		}

		return
	}

	addEtikettToZettel := func(hString string, e etikett.Etikett) (err error) {
		var z zettel_named.Zettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		z.Stored.Zettel.Etiketten.Add(e)
		toUpdate[z.Hinweis.String()] = z

		errors.PrintErrf("Added etikett '%s' to zettel '%s'", e, z.Hinweis)

		return
	}

	removeEtikettFromZettel := func(hString string, e etikett.Etikett) (err error) {
		var z zettel_named.Zettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		z.Stored.Zettel.Etiketten.RemovePrefixes(e)
		toUpdate[z.Hinweis.String()] = z

		errors.PrintErrf("Removed etikett '%s' from zettel '%s'", e, z.Hinweis)

		return
	}

	for _, c := range cs.Removed {
		var e etikett.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = removeEtikettFromZettel(c.Key, e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}
	for _, c := range cs.Added {
		var e etikett.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = addEtikettToZettel(c.Key, e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for bez, etts := range cs.New {
		z := zettel.Zettel{
			Etiketten: etts,
		}

		if err = z.Bezeichnung.Set(bez); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = z.Typ.Set("md"); err != nil {
			err = errors.Wrap(err)
			return
		}

		if c.Umwelt.Konfig.DryRun {
			errors.PrintOutf("[%s] (would create)", z.Bezeichnung)
			continue
		}

		var tz zettel_transacted.Zettel

		if tz, err = store.StoreObjekten().Create(z); err != nil {
			err = errors.Errorf("failed to create zettel: %s", err)
			return
		}

		errors.PrintOutf("%s (created)", tz.Named)
	}

	for _, z := range toUpdate {
		if c.Umwelt.Konfig.DryRun {
			errors.PrintOutf("[%s] (would update)", z.Hinweis)
			continue
		}

		if _, err = store.StoreObjekten().Update(z.Hinweis, z.Stored.Zettel); err != nil {
			errors.PrintErrf("failed to update zettel: %s", err)
		}
	}

	return
}
