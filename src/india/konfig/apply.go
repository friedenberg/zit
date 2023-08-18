package konfig

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/typ"
)

func (k compiled) ApplyToMetadatei(
	ml metadatei.MetadateiLike,
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
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

	var ta *typ.Akte

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
