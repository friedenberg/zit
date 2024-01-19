package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/expansion"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/etiketten_path"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (k *Compiled) ApplyToSku(
	sk *sku.Transacted,
) (err error) {
	mp := &sk.Metadatei

	mp.Verzeichnisse.SetExpandedEtiketten(kennung.ExpandMany[kennung.Etikett](
		mp.GetEtiketten(),
		expansion.ExpanderRight,
	))

	g := gattung.Must(sk.GetGattung())
	isEtikett := g == gattung.Etikett

	if g.HasParents() {
		k.SetHasChanges(true)
	}

	var etikett kennung.Etikett

	if isEtikett {
		err = etikett.Set(sk.Kennung.String())

		if err != nil {
			return
		}
	}

	if isEtikett {
		kennung.ExpandOne[kennung.Etikett](
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
		p1.Add(e)

		impl := k.GetImplicitEtiketten(e)

		if impl.Len() == 0 {
			sk.Metadatei.Verzeichnisse.AddPath(p1)
			return
		}

		if err = impl.EachPtr(
			iter.MakeChain(
				ie.AddPtr,
				func(e1 *kennung.Etikett) (err error) {
					p2 := p1.Copy()
					p2.Add(e1)
					sk.Metadatei.Verzeichnisse.AddPath(p2)
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
		iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
			k.EtikettenHidden,
			func(e *kennung.Etikett) bool {
				return ees.ContainsKey(ees.KeyPtr(e))
			},
		)

	if isHiddenEtikett {
		sk.SetArchiviert(true)
		return
	}

	checkFunc := func(e *kennung.Etikett) bool {
		return k.EtikettenHidden.ContainsKey(k.EtikettenHidden.KeyPtr(e))
	}

	if iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
		mp.GetEtiketten(),
		checkFunc,
	) {
		sk.SetArchiviert(true)
		return
	}

	if iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
		mp.Verzeichnisse.GetExpandedEtiketten(),
		checkFunc,
	) {
		sk.SetArchiviert(true)
		return
	}

	if iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
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

	if !tk.HasValue() {
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
