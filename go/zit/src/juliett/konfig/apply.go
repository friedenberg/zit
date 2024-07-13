package konfig

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/etiketten_path"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO
func (k *Compiled) ApplySchlummerndAndRealizeEtiketten(
	sk *sku.Transacted,
) (err error) {
	ui.Log().Print("applying konfig to:", sk)
	mp := &sk.Metadatei

	mp.Verzeichnisse.SetExpandedEtiketten(kennung.ExpandMany(
		mp.GetEtiketten(),
		expansion.ExpanderRight,
	))

	g := gattung.Must(sk.GetGattung())
	isEtikett := g == gattung.Etikett

	// if g.HasParents() {
	// 	k.SetHasChanges(fmt.Sprintf("adding etikett with parents: %s", sk))
	// }

	var etikett kennung.Etikett

	// TODO better solution for "realizing" etiketten against Konfig.
	// Specifically, making this less fragile and dependent on remembering to do
	// ApplyToSku for each Sku. Maybe a factory?
	mp.Verzeichnisse.Etiketten.Reset()
	mp.GetEtiketten().Each(mp.Verzeichnisse.Etiketten.AddEtikettOld)

	if isEtikett {
		ks := sk.Kennung.String()

		if err = etikett.Set(ks); err != nil {
			return
		}

		sk.Metadatei.Verzeichnisse.Etiketten.AddSelf(catgut.MakeFromString(ks))

		kennung.ExpandOne(
			&etikett,
			expansion.ExpanderRight,
		).EachPtr(
			mp.Verzeichnisse.GetExpandedEtikettenMutable().AddPtr,
		)
	}

	if err = k.addSuperEtiketten(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = k.addImplicitEtiketten(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.SetSchlummernd(k.schlummernd.ContainsSku(sk))

	return
}

func (k *Compiled) addSuperEtiketten(
	sk *sku.Transacted,
) (err error) {
	g := sk.GetGattung()

	var expanded []string
	var ks string

	switch g {
	case gattung.Etikett, gattung.Typ, gattung.Kasten:
		ks = sk.Kennung.String()

		expansion.ExpanderRight.Expand(
			func(v string) (err error) {
				expanded = append(expanded, v)
				return
			},
			ks,
		)

	default:
		return
	}

	for _, ex := range expanded {
		if ex == ks || ex == "" {
			continue
		}

		var ek *sku.Transacted

		if ek, err = k.getEtikettOrKastenOrTyp(ex); err != nil {
			err = errors.Wrapf(err, "Expanded: %q", ex)
			return
		}

		if ek == nil {
			// this is ok because currently, konfig is applied twice. However, this
			// is fragile as the order in which this method is called is
			// non-deterministic and the `GetEtikett` call may request an Etikett we
			// have not processed yet
			continue
		}

		if ek.Metadatei.Verzeichnisse.Etiketten.Paths.Len() <= 1 {
			ui.Log().Print(ks, ex, ek.Metadatei.Verzeichnisse.Etiketten)
			continue
		}

		prefix := catgut.MakeFromString(ex)

		a := &sk.Metadatei.Verzeichnisse.Etiketten
		b := &ek.Metadatei.Verzeichnisse.Etiketten

		ui.Log().Print("a", a)
		ui.Log().Print("b", b)

		ui.Log().Print("prefix", prefix)

		// ui.Log().Print("before", sk.GetKennung(), ex, prefix, a, b)
		// defer ui.Log().Print("after ", sk.GetKennung(), ex, prefix, a, b)

		if err = a.AddSuperFrom(b, prefix); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("a after", a)
	}

	return
}

func (k *Compiled) addImplicitEtiketten(
	sk *sku.Transacted,
) (err error) {
	mp := &sk.Metadatei
	ie := kennung.MakeEtikettMutableSet()

	addImpEts := func(e *kennung.Etikett) (err error) {
		p1 := etiketten_path.MakePathWithType()
		p1.Type = etiketten_path.TypeIndirect
		p1.Add(catgut.MakeFromString(e.String()))

		impl := k.getImplicitEtiketten(e)

		if impl.Len() == 0 {
			sk.Metadatei.Verzeichnisse.Etiketten.AddPathWithType(p1)
			return
		}

		if err = impl.EachPtr(
			iter.MakeChain(
				ie.AddPtr,
				func(e1 *kennung.Etikett) (err error) {
					p2 := p1.Clone()
					p2.Add(catgut.MakeFromString(e1.String()))
					sk.Metadatei.Verzeichnisse.Etiketten.AddPathWithType(p2)
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

	typKonfig := k.getApproximatedTyp(mp.GetTyp()).ApproximatedOrActual()

	if typKonfig != nil {
		typKonfig.GetEtiketten().EachPtr(ie.AddPtr)
		typKonfig.GetEtiketten().EachPtr(addImpEts)
	}

	mp.Verzeichnisse.SetImplicitEtiketten(ie)

	return
}

func (k compiled) ApplyToNewMetadatei(
	ml metadatei.MetadateiLike,
	tagp interfaces.BlobGetterPutter[*typ_akte.V0],
) (err error) {
	// m := ml.GetMetadatei()

	// normalized := kennung.WithRemovedCommonPrefixes(m.GetEtiketten())
	// m.SetEtiketten(normalized)

	return
}
