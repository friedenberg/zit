package konfig

import (
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func init() {
	collections_ptr.RegisterGobValue[ketikett, *ketikett](nil)
}

type implicitEtikettenMap map[kennung.Etikett]kennung.EtikettMutableSet

func (iem implicitEtikettenMap) Contains(to, imp kennung.Etikett) bool {
	s, ok := iem[to]

	if !ok || s == nil {
		return false
	}

	if !s.Contains(imp) {
		return false
	}

	return true
}

func (iem implicitEtikettenMap) Set(to, imp kennung.Etikett) (err error) {
	s, ok := iem[to]

	if !ok {
		s = kennung.MakeEtikettMutableSet()
		iem[to] = s
	}

	return s.Add(imp)
}

type ketikett struct {
	Transacted        sku.Transacted
	ImplicitEtiketten kennung.EtikettMutableSet
}

func (a ketikett) Less(b ketikett) bool {
	return a.Transacted.Less(b.Transacted)
}

func (a ketikett) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a ketikett) Equals(b ketikett) bool {
	if !a.Transacted.Equals(b.Transacted) {
		return false
	}

	if !iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](
		a.ImplicitEtiketten,
		b.ImplicitEtiketten,
	) {
		return false
	}

	return true
}

func (e *ketikett) Set(v string) (err error) {
	return (&e.Transacted.Kennung).Set(v)
}

func (e ketikett) String() string {
	return e.Transacted.GetKennungLike().String()
}

func (k compiled) EachEtikett(
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return k.Etiketten.Each(
		func(ek ketikett) (err error) {
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

func (k *compiled) AddEtikett2(
	b1 *sku.Transacted,
) (err error) {
	e := sku.GetTransactedPool().Get()
	*e = *b1

	return k.AddEtikett(e)
}

func (k *compiled) AddEtikett(
	b1 *sku.Transacted,
) (err error) {
	k.lock.Lock()
	defer k.lock.Unlock()
	k.hasChanges = true

	b := ketikett{
		Transacted: *b1,
	}

	if err = iter.AddOrReplaceIfGreater[ketikett](k.Etiketten, b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *compiled) applyExpandedEtikett(ct *sku.Transacted) {
}

func (kc compiled) GetEtikett(
	k kennung.Etikett,
) (ct *sku.Transacted, ok bool) {
	expandedActual := kc.GetSortedEtikettenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
		ok = true
	}

	return
}

func (c compiled) GetSortedEtikettenExpanded(
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
			ct, ok := c.Etiketten.GetPtr(v.String())

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

func (c compiled) GetImplicitEtiketten(
	e *kennung.Etikett,
) kennung.EtikettSet {
	s, ok := c.ImplicitEtiketten[*e]

	if !ok || s == nil {
		return kennung.MakeEtikettSet()
	}

	return s
}