package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/lima/changes"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommitOrganizeFile struct {
	*umwelt.Umwelt
}

type CommitOrganizeFileResults struct{}

func (c CommitOrganizeFile) Run(
	a, b *organize_text.Text,
) (results CommitOrganizeFileResults, err error) {
	store := c.StoreObjekten()

	var cs changes.Changes

	if cs, err = changes.ChangesFrom(
		a,
		b,
		store.GetAbbrStore().Hinweis().ExpandString,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sameTyp := a.Metadatei.Typ.Equals(b.Metadatei.Typ)

	toUpdate := collections.MakeMutableSetStringer[sku.WithKennungInterface]()

	_ = func(hString string) (z sku.WithKennungInterface, err error) {
		var h kennung.Hinweis

		if h, err = store.GetAbbrStore().Hinweis().ExpandString(
			hString,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		var ok bool

		if z, ok = toUpdate.Get(h.String()); !ok {
			var tz *zettel.Transacted

			if tz, err = store.Zettel().ReadOne(&h); err != nil {
				err = errors.Wrapf(
					err,
					"{Hinweis String: '%s'}: {Hinweis: '%s'}",
					hString,
					h,
				)
				return
			}

			z.Kennung = &tz.Sku.Kennung
			z.Metadatei = tz.GetMetadatei()
		}

		return
	}

	ms := c.Umwelt.MakeMetaIdSetWithoutExcludedHidden(
		nil,
		gattungen.MakeSet(gattung.TrueGattung()...),
	)
	errors.TodoP1("create query without syntax")
	ms.Set(":z,e,t")

	if err = store.Query(
		ms,
		func(tl objekte.TransactedLikePtr) (err error) {
			var change changes.Change
			ok := false
			k := tl.GetKennungPtr()

			if change, ok = cs.GetExisting().Get(kennung.FormattedString(k)); !ok {
				return
			}

			if sameTyp && change.IsEmpty() {
				return
			}

			m := tl.GetMetadatei()
			mes := m.GetEtiketten().MutableClone()
			change.GetRemoved().Each(mes.Del)
			change.GetAdded().Each(mes.Add)
			m.Etiketten = mes.ImmutableClone()

			if !sameTyp {
				m.Typ = b.Metadatei.Typ
			}

			mwk := sku.WithKennungInterface{
				Kennung:   k.KennungPtrClone(),
				Metadatei: m,
			}

			toUpdate.Add(mwk)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = cs.GetAdded().Each(
		func(change changes.Change) (err error) {
			bez := change.Key

			m := sku.Metadatei{
				Etiketten: change.GetAdded().ImmutableClone(),
				Typ:       b.Metadatei.Typ,
			}

			if err = m.Bezeichnung.Set(bez); err != nil {
				err = errors.Wrap(err)
				return
			}

			if kennung.IsEmpty(m.GetTyp()) {
				m.Typ = c.Konfig().Akte.DefaultTyp
			}

			if c.Konfig().DryRun {
				errors.Out().Printf("[%s] (would create)", m.Bezeichnung)
				return
			}

			if _, err = store.Zettel().Create(m); err != nil {
				err = errors.Errorf("failed to create zettel: %s", err)
				return
			}
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if toUpdate.Len() == 0 && cs.GetAdded().Len() == 0 {
		errors.Err().Print("no changes")
		return
	}

	if err = store.UpdateManyMetadatei(toUpdate); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
