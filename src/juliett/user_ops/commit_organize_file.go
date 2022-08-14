package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/logz"
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/india/changes"
	"github.com/friedenberg/zit/src/india/store_with_lock"
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

	changes := changes.ChangesFrom(a, b)

	logz.Printf("%#v", changes)

	if len(changes.Added) == 0 && len(changes.Removed) == 0 && len(changes.New) == 0 {
		stdprinter.Err("no changes")
		return
	}

	toUpdate := make(map[string]stored_zettel.Named)

	addOrGetToZettelToUpdate := func(hString string) (z stored_zettel.Named, err error) {
		var h hinweis.Hinweis

		if h, err = hinweis.MakeBlindHinweis(hString); err != nil {
			err = errors.Error(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			var tz stored_zettel.Transacted

			if tz, err = store.Zettels().Read(h); err != nil {
				err = errors.Error(err)
				return
			}

			z = tz.Named
		}

		return
	}

	addEtikettToZettel := func(hString string, e etikett.Etikett) (err error) {
		var z stored_zettel.Named

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Error(err)
			return
		}

		z.Zettel.Etiketten.Add(e)
		toUpdate[z.Hinweis.String()] = z

		stdprinter.Errf("Added etikett '%s' to zettel '%s'\n", e, z.Hinweis)

		return
	}

	removeEtikettFromZettel := func(hString string, e etikett.Etikett) (err error) {
		var z stored_zettel.Named

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Error(err)
			return
		}

		z.Zettel.Etiketten.RemovePrefixes(e)
		toUpdate[z.Hinweis.String()] = z

		stdprinter.Errf("Removed etikett '%s' from zettel '%s'\n", e, z.Hinweis)

		return
	}

	for _, c := range changes.Removed {
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
	for _, c := range changes.Added {
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

	for bez, etts := range changes.New {
		z := zettel.Zettel{
			Etiketten: etts,
		}

		if err = z.Bezeichnung.Set(bez); err != nil {
			err = errors.Error(err)
			return
		}

		if err = z.AkteExt.Set("md"); err != nil {
			err = errors.Error(err)
			return
		}

		if c.Umwelt.Konfig.DryRun {
			stdprinter.Outf("[%s] (would create)\n", z.Bezeichnung)
			continue
		}

		var tz stored_zettel.Transacted

		if tz, err = store.Zettels().Create(z); err != nil {
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

		if _, err = store.Zettels().Update(z.Hinweis, z.Zettel); err != nil {
			stdprinter.Errf("failed to update zettel: %s", err)
		}
	}

	return
}
