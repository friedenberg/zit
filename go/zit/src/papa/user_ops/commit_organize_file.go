package user_ops

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/lima/changes"
	"code.linenisgreat.com/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommitOrganizeFile struct {
	*umwelt.Umwelt
	OutputJSON bool
}

type CommitOrganizeFileResults struct{}

func (c CommitOrganizeFile) ApplyToText(
	u *umwelt.Umwelt,
	t *organize_text.Text,
) (err error) {
	if u.Konfig().PrintOptions.PrintEtikettenAlways {
		return
	}

	if err = t.Transacted.Each(
		func(sk *sku.Transacted) (err error) {
			if sk.Metadatei.Bezeichnung.IsEmpty() {
				return
			}

			sk.Metadatei.GetEtikettenMutable().Reset()

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op CommitOrganizeFile) Run(
	u *umwelt.Umwelt,
	a, b *organize_text.Text,
) (results CommitOrganizeFileResults, err error) {
	store := op.StoreObjekten()

	if err = op.ApplyToText(u, a); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	ms := op.MakeMetaIdSetWithoutExcludedHidden(
		gattungen.MakeSet(gattung.TrueGattung()...),
	)
	errors.TodoP1("create query without syntax")
	ms.Set(":z,e,t")

	if err = store.QueryWithCwd(
		ms,
		func(tl *sku.Transacted) (err error) {
			var change *changes.Change
			ok := false
			sk := sku.GetTransactedPool().Get()

			if err = sk.SetFromSkuLike(tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			k := kennung.FormattedString(&sk.Kennung)

			if change, ok = cs.GetExisting().Get(k); !ok {
				return
			}

			bezChange, didBezChange := cs.GetModified().Get(k)

			if sameTyp && change.IsEmpty() && !didBezChange {
				return
			}

			m := sk.GetMetadatei()
			change.GetRemoved().EachPtr(m.GetEtikettenMutable().DelPtr)
			change.GetAdded().EachPtr(m.AddEtikettPtr)

			if didBezChange {
				m.Bezeichnung = bezChange.Bezeichnung
			}

			if !sameTyp {
				m.Typ = b.Metadatei.Typ
			}

			l.Lock()
			defer l.Unlock()

			if err = toUpdate.Add(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = cs.GetAddedNamed().Each(
		func(change *changes.Change) (err error) {
			var k kennung.Kennung2

			if err = k.Set(change.Key); err != nil {
				err = errors.Wrap(err)
				return
			}

			m := &metadatei.Metadatei{
				Typ: b.Metadatei.Typ,
			}

			m.SetEtiketten(change.GetAdded())

			// TODO-P2 add support for setting bez
			// if err = m.Bezeichnung.Set(bez); err != nil {
			// 	err = errors.Wrap(err)
			// 	return
			// }

			// TODO-P2 determine appropriate typ for named
			// if kennung.IsEmpty(m.GetTyp()) {
			// 	m.Typ = c.Konfig().Defaults.Typ
			// }

			if op.Konfig().DryRun {
				errors.Out().Printf("[%s] (would create)", k)
				return
			}

			if _, err = store.CreateOrUpdate(m, &k); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = cs.GetAddedUnnamed().Each(
		func(change *changes.Change) (err error) {
			bez := change.Key

			m := &metadatei.Metadatei{
				Typ: b.Metadatei.Typ,
			}

			m.SetEtiketten(change.GetAdded())

			if err = m.Bezeichnung.Set(bez); err != nil {
				err = errors.Wrap(err)
				return
			}

			if kennung.IsEmpty(m.GetTyp()) {
				m.Typ = op.Konfig().Defaults.Typ
			}

			if op.Konfig().DryRun {
				errors.Out().Printf("[%s] (would create)", m.Bezeichnung)
				return
			}

			if _, err = store.Create(m); err != nil {
				err = errors.Errorf("failed to create zettel: %s", err)
				return
			}
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if toUpdate.Len() == 0 && cs.GetAddedUnnamed().Len() == 0 && cs.GetAddedNamed().Len() == 0 {
		errors.Err().Print("no changes")
		return
	}

	if err = store.UpdateManyMetadatei(toUpdate); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.OutputJSONIfNecessary(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op CommitOrganizeFile) OutputJSONIfNecessary() (err error) {
	if !op.OutputJSON {
		return
	}

	// if err = createOrganizeFileResults.EachPtr(
	// 	func(sk *sku.Transacted) (err error) {
	// 		var j sku_fmt.Json

	// 		if err = j.FromTransacted(sk, u.Standort()); err != nil {
	// 			err = errors.Wrap(err)
	// 			return
	// 		}

	// 		transacted = append(transacted, j)

	// 		return
	// 	},
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	// w := bufio.NewWriter(u.Out())
	// defer errors.DeferredFlusher(w)
	// enc := json.NewEncoder(w)

	return
}
