package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/srx/bravo/expansion"
)

func (k compiled) ApplyToMetadatei(
	ml metadatei.MetadateiLike,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	mp := ml.GetMetadateiPtr()
	mp.Verzeichnisse.SetExpandedEtiketten(kennung.ExpandMany[kennung.Etikett](
		mp.GetEtiketten(),
		expansion.ExpanderRight,
	))

	ie := kennung.MakeEtikettMutableSet()

	addImpEts := func(e *kennung.Etikett) (err error) {
		impl := k.GetImplicitEtiketten(e)

		if err = impl.EachPtr(ie.AddPtr); err != nil {
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

	checkFunc := func(e *kennung.Etikett) bool {
		return k.EtikettenHidden.ContainsKey(k.EtikettenHidden.KeyPtr(e))
	}

	mp.Verzeichnisse.Archiviert.SetBool(
		iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
			mp.GetEtiketten(),
			checkFunc,
		) ||
			iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
				mp.Verzeichnisse.GetExpandedEtiketten(),
				checkFunc,
			) ||
			iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
				mp.Verzeichnisse.GetImplicitEtiketten(),
				checkFunc,
			),
	)

	return
}

func (k compiled) ApplyToNewMetadatei(
	ml metadatei.MetadateiLike,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	m := ml.GetMetadatei()

	defer func() {
		if err == nil {
			ml.SetMetadatei(m)
		}
	}()

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