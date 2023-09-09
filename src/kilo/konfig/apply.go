package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

func (k compiled) ApplyToMetadatei(
	ml metadatei.MetadateiLike,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	mp := ml.GetMetadateiPtr()
	mp.Verzeichnisse.ExpandedEtiketten = kennung.ExpandMany[kennung.Etikett](
		mp.GetEtiketten(),
		kennung.ExpanderRight,
	)

	ie := kennung.MakeEtikettMutableSet()

	mp.GetEtiketten().EachPtr(
		func(e *kennung.Etikett) (err error) {
			impl := k.GetImplicitEtiketten(e)

			if err = impl.EachPtr(ie.AddPtr); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	mp.Verzeichnisse.ImplicitEtiketten = ie

	checkFunc := func(e *kennung.Etikett) bool {
		return k.EtikettenHidden.ContainsKey(k.EtikettenHidden.KeyPtr(e))
	}

	mp.Verzeichnisse.Archiviert.SetBool(
		iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
			mp.Etiketten,
			checkFunc,
		) ||
			iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
				mp.Verzeichnisse.ExpandedEtiketten,
				checkFunc,
			) ||
			iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
				mp.Verzeichnisse.ImplicitEtiketten,
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
	normalized := kennung.WithRemovedCommonPrefixes(m.Etiketten)
	m.Etiketten = normalized

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
