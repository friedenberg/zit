package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/stdprinter"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/golf/organize_text"
	"github.com/friedenberg/zit/india/changes"
	"github.com/friedenberg/zit/india/store_with_lock"
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

	if len(changes.Added) == 0 && len(changes.Removed) == 0 && len(changes.New) == 0 {
		stdprinter.Err("no changes")
		return
	}

	toUpdate := make(map[string]_NamedZettel)

	addOrGetToZettelToUpdate := func(hString string) (z _NamedZettel, err error) {
		var h hinweis.Hinweis

		if h, err = hinweis.MakeBlindHinweis(hString); err != nil {
			err = _Error(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			if z, err = store.Zettels().Read(h); err != nil {
				err = _Error(err)
				return
			}
		}

		return
	}

	addEtikettToZettel := func(hString string, e etikett.Etikett) (err error) {
		var z _NamedZettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = _Error(err)
			return
		}

		z.Zettel.Etiketten.Add(e)
		toUpdate[z.Hinweis.String()] = z

		stdprinter.Errf("Added etikett '%s' to zettel '%s'\n", e, z.Hinweis)

		return
	}

	removeEtikettFromZettel := func(hString string, e etikett.Etikett) (err error) {
		var z _NamedZettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = _Error(err)
			return
		}

		z.Zettel.Etiketten.RemovePrefixes(e)
		toUpdate[z.Hinweis.String()] = z

		stdprinter.Errf("Removed etikett '%s' from zettel '%s'\n", e, z.Hinweis)

		return
	}

	for _, c := range changes.Added {
		var e etikett.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = _Error(err)
			return
		}

		if err = addEtikettToZettel(c.Key, e); err != nil {
			err = _Error(err)
			return
		}
	}

	for _, c := range changes.Removed {
		var e etikett.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = _Error(err)
			return
		}

		if err = removeEtikettFromZettel(c.Key, e); err != nil {
			err = _Error(err)
			return
		}
	}

	for bez, etts := range changes.New {
		z := _Zettel{
			Etiketten: etts,
		}

		if err = z.Bezeichnung.Set(bez); err != nil {
			err = _Error(err)
			return
		}

		if err = z.AkteExt.Set("md"); err != nil {
			err = _Error(err)
			return
		}

		if c.Umwelt.Konfig.DryRun {
			stdprinter.Outf("[%s] (would create)\n", z.Bezeichnung)
			continue
		}

		var named _NamedZettel

		if named, err = store.Zettels().Create(z); err != nil {
			stdprinter.Errf("failed to create zettel: %s", err)
		}

		stdprinter.Outf("[%s %s] (created)\n", named.Hinweis, named.Sha)
	}

	for _, z := range toUpdate {
		if c.Umwelt.Konfig.DryRun {
			_Outf("[%s] (would update)\n", z.Hinweis)
			continue
		}

		if _, err = store.Zettels().Update(z); err != nil {
			stdprinter.Errf("failed to update zettel: %s", err)
		}
	}

	return
}
