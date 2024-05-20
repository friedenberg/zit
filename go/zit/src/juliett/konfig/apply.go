package konfig

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/etiketten_path"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func (k *Compiled) ApplyToSku(
	sk *sku.Transacted,
) (err error) {
	mp := &sk.Metadatei

	mp.Verzeichnisse.SetExpandedEtiketten(kennung.ExpandMany(
		mp.GetEtiketten(),
		expansion.ExpanderRight,
	))

	g := gattung.Must(sk.GetGattung())
	isEtikett := g == gattung.Etikett

	if g.HasParents() {
		k.SetHasChanges()
	}

	var etikett kennung.Etikett

	if isEtikett {
		err = etikett.Set(sk.Kennung.String())
		if err != nil {
			return
		}
	}

	if isEtikett {
		kennung.ExpandOne(
			&etikett,
			expansion.ExpanderRight,
		).EachPtr(
			mp.Verzeichnisse.GetExpandedEtikettenMutable().AddPtr,
		)
	}

	if err = k.addImplicitEtiketten(sk, nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = k.setArchiviert(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *Compiled) addImplicitEtiketten(
	sk *sku.Transacted,
	p *etiketten_path.Path,
) (err error) {
	mp := &sk.Metadatei
	ie := kennung.MakeEtikettMutableSet()

	addImpEts := func(e *kennung.Etikett) (err error) {
		p1 := p.Copy()
		p1.Add(catgut.MakeFromString(e.String()))

		impl := k.GetImplicitEtiketten(e)

		if impl.Len() == 0 {
			sk.Metadatei.Verzeichnisse.Etiketten.AddPath(p1)
			return
		}

		if err = impl.EachPtr(
			iter.MakeChain(
				ie.AddPtr,
				func(e1 *kennung.Etikett) (err error) {
					p2 := p1.Copy()
					p2.Add(catgut.MakeFromString(e1.String()))
					sk.Metadatei.Verzeichnisse.Etiketten.AddPath(p2)
					return
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	mp.GetEtiketten().EachPtr(addImpEts)

	typKonfig := k.GetApproximatedTyp(mp.GetTyp()).ApproximatedOrActual()

	if typKonfig != nil {
		typKonfig.GetEtiketten().EachPtr(ie.AddPtr)
		typKonfig.GetEtiketten().EachPtr(addImpEts)
	}

	mp.Verzeichnisse.SetImplicitEtiketten(ie)

	return
}

func (k *Compiled) setArchiviert(
	sk *sku.Transacted,
) (err error) {
	sk.SetArchiviert(false)

	mp := &sk.Metadatei

	g := gattung.Must(sk.GetGattung())
	isEtikett := g == gattung.Etikett

	ees := mp.Verzeichnisse.GetExpandedEtiketten()

	isHiddenEtikett := isEtikett &&
		iter.CheckAnyPtr(
			k.EtikettenHidden,
			func(e *kennung.Etikett) bool {
				ok := ees.ContainsKey(ees.KeyPtr(e))
				return ok
			},
		)

	if isHiddenEtikett {
		sk.SetArchiviert(true)
		return
	}

	checkFunc := func(e *kennung.Etikett) bool {
		ok := k.EtikettenHidden.ContainsKey(k.EtikettenHidden.KeyPtr(e))
		return ok
	}

	if iter.CheckAnyPtr(
		mp.GetEtiketten(),
		checkFunc,
	) {
		sk.SetArchiviert(true)
		return
	}

	if iter.CheckAnyPtr(
		mp.Verzeichnisse.GetExpandedEtiketten(),
		checkFunc,
	) {
		sk.SetArchiviert(true)
		return
	}

	if iter.CheckAnyPtr(
		mp.Verzeichnisse.GetImplicitEtiketten(),
		checkFunc,
	) {
		sk.SetArchiviert(true)
		return
	}

	return
}

func (k compiled) ApplyToNewMetadatei(
	ml metadatei.MetadateiLike,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	m := ml.GetMetadatei()

	t := m.GetTyp()
	normalized := kennung.WithRemovedCommonPrefixes(m.GetEtiketten())
	m.SetEtiketten(normalized)

	tk := k.GetApproximatedTyp(t)

	if !tk.HasValue {
		return
	}

	toa := tk.ApproximatedOrActual()

	if toa == nil {
		return
	}

	var ta *typ_akte.V0

	if ta, err = tagp.GetAkte(toa.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	for e, r := range ta.EtikettenRules {
		var e1 kennung.Etikett

		if e1, err = kennung.MakeEtikett(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = m.ApplyGoldenChild(e1, r.GoldenChild); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
