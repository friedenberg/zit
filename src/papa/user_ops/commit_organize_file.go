package user_ops

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/organize_text"
	"github.com/friedenberg/zit/src/lima/changes"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

	l := sync.Mutex{}
	toUpdate := sku.MakeTransactedMutableSet()

	ms := c.Umwelt.MakeMetaIdSetWithoutExcludedHidden(
		gattungen.MakeSet(gattung.TrueGattung()...),
	)
	errors.TodoP1("create query without syntax")
	ms.Set(":z,e,t")

	if err = store.Query(
		ms,
		func(tl *sku.Transacted) (err error) {
			var change changes.Change
			ok := false
			sk := sku.GetTransactedPool().Get()

			if err = sk.SetFromSkuLike(tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			k := sk.GetKennungLike()

			if change, ok = cs.GetExisting().Get(kennung.FormattedString(k)); !ok {
				return
			}

			if sameTyp && change.IsEmpty() {
				return
			}

			m := sk.GetMetadatei()
			change.GetRemoved().EachPtr(m.GetEtikettenMutable().DelPtr)
			change.GetAdded().EachPtr(m.AddEtikettPtr)

			if !sameTyp {
				m.Typ = b.Metadatei.Typ
			}

			sk.SetMetadatei(m)

			l.Lock()
			defer l.Unlock()

			if err = toUpdate.AddPtr(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = cs.GetAdded().Each(
		func(change changes.Change) (err error) {
			bez := change.Key

			m := metadatei.Metadatei{
				Typ: b.Metadatei.Typ,
			}

			m.SetEtiketten(change.GetAdded())

			if err = m.Bezeichnung.Set(bez); err != nil {
				err = errors.Wrap(err)
				return
			}

			if kennung.IsEmpty(m.GetTyp()) {
				m.Typ = c.Konfig().Akte.Defaults.Typ
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