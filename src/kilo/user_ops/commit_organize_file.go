package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/organize_text"
	"github.com/friedenberg/zit/src/india/changes"
	"github.com/friedenberg/zit/src/juliett/zettel_printer"
)

type CommitOrganizeFile struct {
	*zettel_printer.Printer
}

type CommitOrganizeFileResults struct {
}

func (c CommitOrganizeFile) Run(a, b organize_text.Text) (results CommitOrganizeFileResults, err error) {
	store := c.Store

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

			if tz, err = store.Read(h); err != nil {
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

		if tz, err = store.Create(z); err != nil {
			err = errors.Errorf("failed to create zettel: %s", err)
			return
		}

		c.ZettelTransacted(tz).Print()

		if !c.IsEmpty() {
			err = c.Error()
			return
		}
	}

	for _, z := range toUpdate {
		if c.Umwelt.Konfig.DryRun {
			errors.PrintOutf("[%s] (would update)", z.Hinweis)
			continue
		}

		if _, err = store.Update(z.Hinweis, z.Stored.Zettel); err != nil {
			errors.PrintErrf("failed to update zettel: %s", err)
		}
	}

	return
}
