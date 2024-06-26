package konfig

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func init() {
	collections_value.RegisterGobValue[*ketikett](nil)
}

type implicitEtikettenMap map[string]kennung.EtikettMutableSet

func (iem implicitEtikettenMap) Contains(to, imp kennung.Etikett) bool {
	s, ok := iem[to.String()]

	if !ok || s == nil {
		return false
	}

	if !s.Contains(imp) {
		return false
	}

	return true
}

func (iem implicitEtikettenMap) Set(to, imp kennung.Etikett) (err error) {
	s, ok := iem[to.String()]

	if !ok {
		s = kennung.MakeEtikettMutableSet()
		iem[to.String()] = s
	}

	return s.Add(imp)
}

type ketikett struct {
	Transacted sku.Transacted
	Computed   bool
}

func (a *ketikett) Less(b *ketikett) bool {
	return sku.TransactedLessor.Less(&a.Transacted, &b.Transacted)
}

func (a *ketikett) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a *ketikett) Equals(b *ketikett) bool {
	if !a.Transacted.Equals(&b.Transacted) {
		return false
	}

	if !iter.SetEqualsPtr(
		a.Transacted.Metadatei.Verzeichnisse.GetImplicitEtiketten(),
		b.Transacted.Metadatei.Verzeichnisse.GetImplicitEtiketten(),
	) {
		return false
	}

	return true
}

func (e *ketikett) Set(v string) (err error) {
	return (&e.Transacted.Kennung).Set(v)
}

func (e *ketikett) String() string {
	return e.Transacted.GetKennung().String()
}

func (k *compiled) EachEtikett(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return k.Etiketten.Each(
		func(ek *ketikett) (err error) {
			return f(&ek.Transacted)
		},
	)
}

func (k *compiled) AccumulateImplicitEtiketten(
	e kennung.Etikett,
) (err error) {
	ek, ok := k.Etiketten.Get(e.String())

	if !ok {
		return
	}

	ees := kennung.MakeEtikettMutableSet()

	kennung.ExpandOne(&e, expansion.ExpanderRight).EachPtr(
		ees.AddPtr,
	)

	if err = ees.Each(
		func(e1 kennung.Etikett) (err error) {
			if e1.Equals(e) {
				return
			}

			if err = k.AccumulateImplicitEtiketten(e1); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = k.getImplicitEtiketten(&e1).Each(
				func(e2 kennung.Etikett) (err error) {
					return k.ImplicitEtiketten.Set(e, e2)
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = ek.Transacted.Metadatei.GetEtiketten().Each(
		func(e1 kennung.Etikett) (err error) {
			if k.ImplicitEtiketten.Contains(e1, e) {
				return
			}

			if err = k.ImplicitEtiketten.Set(e, e1); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = k.AccumulateImplicitEtiketten(e1); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k *compiled) addEtikett(
	kinder *sku.Transacted,
	mutter *sku.Transacted,
) (didChange bool, err error) {
	k.lock.Lock()
	defer k.lock.Unlock()

	var b ketikett

	if err = b.Transacted.SetFromSkuLike(kinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	if didChange, err = iter.AddOrReplaceIfGreater(k.Etiketten, &b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
