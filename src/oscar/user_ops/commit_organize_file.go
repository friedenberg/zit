package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/juliett/organize_text"
	"github.com/friedenberg/zit/src/kilo/changes"
	"github.com/friedenberg/zit/src/kilo/zettel"
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

	errors.Log().Printf("%#v", cs)

	sameTyp := a.Metadatei.Typ.Equals(&b.Metadatei.Typ)

	if len(cs.Added) == 0 && len(cs.Removed) == 0 && len(cs.New) == 0 && sameTyp {
		errors.Err().Print("no changes")
		return
	}

	type zettelToUpdate struct {
		objekte zettel.Objekte
		kennung hinweis.Hinweis
	}

	toUpdate := make(map[string]zettelToUpdate)

	addOrGetToZettelToUpdate := func(hString string) (z zettelToUpdate, err error) {
		var h hinweis.Hinweis

		if h, err = store.Abbr().ExpandHinweisString(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ok bool

		if z, ok = toUpdate[h.String()]; !ok {
			var tz zettel.Transacted

			if tz, err = store.Zettel().ReadHinweisSchwanzen(h); err != nil {
				err = errors.Wrapf(err, "{Hinweis String: '%s'}: {Hinweis: '%s'}", hString, h)
				return
			}

			z.objekte = tz.Objekte
			z.kennung = tz.Sku.Kennung
		}

		return
	}

	addEtikettToZettel := func(hString string, e kennung.Etikett) (err error) {
		var z zettelToUpdate

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		mes := z.objekte.Etiketten.MutableCopy()
		mes.Add(e)
		z.objekte.Etiketten = mes.Copy()
		toUpdate[z.kennung.String()] = z

		errors.Err().Printf("Added etikett '%s' to zettel '%s'", e, z.kennung)

		return
	}

	removeEtikettFromZettel := func(hString string, e kennung.Etikett) (err error) {
		var z zettelToUpdate

		if z, err = addOrGetToZettelToUpdate(hString); err != nil {
			err = errors.Wrap(err)
			return
		}

		mes := z.objekte.Etiketten.MutableCopy()
		mes.RemovePrefixes(e)
		z.objekte.Etiketten = mes.Copy()

		toUpdate[z.kennung.String()] = z

		errors.Err().Printf("Removed etikett '%s' from zettel '%s'", e, z.kennung)

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
			var z zettelToUpdate

			if z, err = addOrGetToZettelToUpdate(h); err != nil {
				err = errors.Wrap(err)
				return
			}

			z.objekte.Typ = b.Metadatei.Typ

			toUpdate[z.kennung.String()] = z

			errors.Err().Printf("Switched to typ '%s' for zettel '%s'", b.Metadatei.Typ, z.kennung)
		}
	}

	for _, n := range cs.New {
		bez := n.Key
		etts := n.Etiketten

		z := zettel.Objekte{
			Etiketten: etts.Copy(),
			Typ:       b.Metadatei.Typ,
		}

		if err = z.Bezeichnung.Set(bez); err != nil {
			err = errors.Wrap(err)
			return
		}

		if z.Typ.IsEmpty() {
			if err = z.Typ.Set(
				c.Konfig().DefaultTyp.Sku.Kennung.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if c.Konfig().DryRun {
			errors.Out().Printf("[%s] (would create)", z.Bezeichnung)
			continue
		}

		if _, err = store.Zettel().Create(z); err != nil {
			err = errors.Errorf("failed to create zettel: %s", err)
			return
		}
	}

	for _, z := range toUpdate {
		if c.Konfig().DryRun {
			errors.Out().Printf("[%s] (would update)", z.kennung)
			continue
		}

		if _, err = store.Zettel().Update(&z.objekte, &z.kennung); err != nil {
			errors.Err().Printf("failed to update zettel: %s", err)
		}
	}

	return
}
