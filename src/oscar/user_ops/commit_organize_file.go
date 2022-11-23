package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/kilo/changes"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommitOrganizeFile struct {
	*umwelt.Umwelt
}

type CommitOrganizeFileResults struct {
}

func (c CommitOrganizeFile) Run(a, b *organize_text.Text) (results CommitOrganizeFileResults, err error) {
	store := c.StoreObjekten()

	var cs changes.Changes

	if cs, err = changes.ChangesFrom(a, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	errors.Printf("%#v", cs)

	sameTyp := a.Metadatei.Typ.Equals(b.Metadatei.Typ)

	if len(cs.Added) == 0 && len(cs.Removed) == 0 && len(cs.New) == 0 && sameTyp {
		errors.PrintErr("no changes")
		return
	}

	toUpdate := make(map[string]zettel_named.Zettel)

	addOrGetToZettelToUpdate := func(hString string) (z zettel_named.Zettel, err error) {
		var h hinweis.Hinweis

		if h, err = store.ExpandHinweisString(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			var tz zettel_transacted.Zettel

			if tz, err = store.ReadHinweisSchwanzen(h); err != nil {
				err = errors.Wrapf(err, "{Hinweis String: '%s'}: {Hinweis: '%s'}", hString, h)
				return
			}

			z = tz.Named
		}

		return
	}

	addEtikettToZettel := func(hString string, e kennung.Etikett) (err error) {
		var z zettel_named.Zettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		mes := z.Stored.Objekte.Etiketten.MutableCopy()
		mes.Add(e)
		z.Stored.Objekte.Etiketten = mes.Copy()
		toUpdate[z.Kennung.String()] = z

		errors.PrintErrf("Added etikett '%s' to zettel '%s'", e, z.Kennung)

		return
	}

	removeEtikettFromZettel := func(hString string, e kennung.Etikett) (err error) {
		var z zettel_named.Zettel

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		mes := z.Stored.Objekte.Etiketten.MutableCopy()
		mes.RemovePrefixes(e)
		z.Stored.Objekte.Etiketten = mes.Copy()

		toUpdate[z.Kennung.String()] = z

		errors.PrintErrf("Removed etikett '%s' from zettel '%s'", e, z.Kennung)

		return
	}

	for _, c := range cs.Removed {
		var e kennung.Etikett

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
		var e kennung.Etikett

		if err = e.Set(c.Etikett); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = addEtikettToZettel(c.Key, e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !sameTyp {
		for _, h := range cs.AllB {
			var z zettel_named.Zettel

			if z, err = addOrGetToZettelToUpdate(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			z.Stored.Objekte.Typ = b.Metadatei.Typ

			toUpdate[z.Kennung.String()] = z

			errors.PrintErrf("Switched to typ '%s' for zettel '%s'", b.Metadatei.Typ, z.Kennung)
		}
	}

	for _, n := range cs.New {
		bez := n.Key
		etts := n.Etiketten

		z := zettel.Zettel{
			Etiketten: etts.Copy(),
			Typ:       b.Metadatei.Typ,
		}

		if err = z.Bezeichnung.Set(bez); err != nil {
			err = errors.Wrap(err)
			return
		}

		if z.Typ.IsEmpty() {
			if err = z.Typ.Set(c.Konfig().DefaultTyp.Name.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if c.Konfig().DryRun {
			errors.PrintOutf("[%s] (would create)", z.Bezeichnung)
			continue
		}

		if _, err = store.Create(z); err != nil {
			err = errors.Errorf("failed to create zettel: %s", err)
			return
		}
	}

	for _, z := range toUpdate {
		if c.Konfig().DryRun {
			errors.PrintOutf("[%s] (would update)", z.Kennung)
			continue
		}

		if _, err = store.Update(&z); err != nil {
			errors.PrintErrf("failed to update zettel: %s", err)
		}
	}

	return
}
