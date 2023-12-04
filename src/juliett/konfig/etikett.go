package konfig

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/bravo/expansion"
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

	if !iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](
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

	kennung.ExpandOne[kennung.Etikett](
		&e,
		expansion.ExpanderRight,
	).EachPtr(
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

			if err = k.GetImplicitEtiketten(&e1).Each(
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

func (k *compiled) AddEtikett(
	b1 *sku.Transacted,
) (err error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.hasChanges = true

	var b ketikett

	if err = b.Transacted.SetFromSkuLike(b1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = iter.AddOrReplaceIfGreater[*ketikett](k.Etiketten, &b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *compiled) GetEtikett(
	k kennung.Etikett,
) (ct *sku.Transacted, ok bool) {
	expandedActual := kc.GetSortedEtikettenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
		ok = true
	}

	return
}

func (c *compiled) GetSortedEtikettenExpanded(
	v string,
) (expandedActual []*sku.Transacted) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)
	sa := iter.MakeFuncSetString[values.String, *values.String](
		expandedMaybe,
	)
	typExpander.Expand(sa, v)
	expandedActual = make([]*sku.Transacted, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			ct, ok := c.Etiketten.Get(v.String())

			if !ok {
				return
			}

			ct1 := sku.GetTransactedPool().Get()

			if err = ct1.SetFromSkuLike(&ct.Transacted); err != nil {
				err = errors.Wrap(err)
				return
			}

			expandedActual = append(expandedActual, ct1)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return len(
			expandedActual[i].GetKennung().String(),
		) > len(
			expandedActual[j].GetKennung().String(),
		)
	})

	return
}

func (c *compiled) GetImplicitEtiketten(
	e *kennung.Etikett,
) kennung.EtikettSet {
	s, ok := c.ImplicitEtiketten[e.String()]

	if !ok || s == nil {
		return kennung.MakeEtikettSet()
	}

	return s
}
