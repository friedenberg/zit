package user_ops

import (
	"bufio"
	"encoding/json"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/lima/organize_text"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type CommitOrganizeFile struct {
	*umwelt.Umwelt
	OutputJSON bool
}

type CommitOrganizeFileResults = organize_text.Changes

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
	original sku.TransactedSet,
) (results CommitOrganizeFileResults, err error) {
	if results, err = op.run(u, a, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = op.OutputJSONIfNecessary(results, original); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (op CommitOrganizeFile) run(
	u *umwelt.Umwelt,
	a, b *organize_text.Text,
) (cs CommitOrganizeFileResults, err error) {
	store := op.StoreObjekten()

	if err = op.ApplyToText(u, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	if cs, err = organize_text.ChangesFrom(
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

	builder := op.MakeMetaIdSetWithoutExcludedHidden(
		kennung.MakeGattung(gattung.TrueGattung()...),
	)

	var ids *query.QueryGroup

	errors.TodoP1("create query without syntax")
	if ids, err = builder.BuildQueryGroup(":z,e,t"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = store.QueryWithCwd(
		ids,
		func(tl *sku.Transacted) (err error) {
			var change *organize_text.Change
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
		func(change *organize_text.Change) (err error) {
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
		func(change *organize_text.Change) (err error) {
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

	return
}

func (op CommitOrganizeFile) OutputJSONIfNecessary(
	c organize_text.Changes,
	original sku.TransactedSet,
) (err error) {
	if !op.OutputJSON {
		return
	}

	_, a, b := c.GetChanges()

	var skus, oldSkus sku.TransactedSet

	if skus, err = b.GetSkus(original); err != nil {
		err = errors.Wrap(err)
		return
	}

	if oldSkus, err = a.GetSkus(original); err != nil {
		err = errors.Wrap(err)
		return
	}

	output := map[string][]sku_fmt.Json{
		"changed": make([]sku_fmt.Json, 0),
		"new":     make([]sku_fmt.Json, 0),
		"removed": make([]sku_fmt.Json, 0),
		"same":    make([]sku_fmt.Json, 0),
	}

	if err = oldSkus.Each(
		func(sk *sku.Transacted) (err error) {
			var j sku_fmt.Json

			if err = j.FromTransacted(sk, op.Standort()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, ok := skus.Get(oldSkus.Key(sk)); ok {
				return
			}

			output["removed"] = append(output["removed"], j)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = skus.Each(
		func(sk *sku.Transacted) (err error) {
			var j sku_fmt.Json

			if err = j.FromTransacted(sk, op.Standort()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if old, ok := oldSkus.Get(oldSkus.Key(sk)); ok {
				if sku.TransactedEqualer.Equals(sk, old) {
					output["same"] = append(output["same"], j)
				} else {
					output["changed"] = append(output["changed"], j)
				}
			} else {
				output["new"] = append(output["new"], j)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	w := bufio.NewWriter(op.Out())
	defer errors.DeferredFlusher(&err, w)
	enc := json.NewEncoder(w)

	if err = enc.Encode(output); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
