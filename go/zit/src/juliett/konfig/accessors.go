package konfig

import (
	"fmt"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (c *compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.FileExtensions.Zettel)
}

func (kc *Compiled) GetAngeboren() schnittstellen.Angeboren {
	return kc.angeboren
}

func (kc *compiled) getTyp(k kennung.Kennung) (ct *sku.Transacted) {
	if k.GetGattung() != gattung.Typ {
		return
	}

	if ct1, ok := kc.Typen.Get(k.String()); ok {
		ct = sku.GetTransactedPool().Get()
		errors.PanicIfError(ct.SetFromSkuLike(ct1))
	}

	return
}

func (kc *compiled) getKasten(k kennung.Kennung) (ct *sku.Transacted) {
	if k.GetGattung() != gattung.Kasten {
		return
	}

	if ct1, ok := kc.Kisten.Get(k.String()); ok {
		ct = sku.GetTransactedPool().Get()
		errors.PanicIfError(ct.SetFromSkuLike(ct1))
	}

	return
}

// Returns the exactly matching Typ, or if it doesn't exist, returns the parent
// Typ or nil. (Parent Typ for `md-gdoc` would be `md`.)
func (kc *compiled) getApproximatedTyp(
	k kennung.Kennung,
) (ct ApproximatedTyp) {
	if k.GetGattung() != gattung.Typ {
		return
	}

	expandedActual := kc.getSortedTypenExpanded(k.String())
	if len(expandedActual) > 0 {
		ct.HasValue = true
		ct.Typ = expandedActual[0]

		if kennung.Equals(ct.Typ.GetKennung(), k) {
			ct.IsActual = true
		}
	}

	return
}

func (kc *compiled) getEtikettOrKastenOrTyp(
	v string,
) (sk *sku.Transacted, err error) {
	var k kennung.Kennung2

	if err = k.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch k.GetGattung() {
	case gattung.Etikett:
		sk, _ = kc.getEtikett(&k)
	case gattung.Kasten:
		sk = kc.getKasten(&k)
	case gattung.Typ:
		sk = kc.getTyp(&k)

	default:
		err = gattung.MakeErrUnsupportedGattung(&k)
		return
	}

	return
}

func (kc *compiled) getEtikett(
	k kennung.Kennung,
) (ct *sku.Transacted, ok bool) {
	if k.GetGattung() != gattung.Etikett {
		return
	}

	expandedActual := kc.getSortedEtikettenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
		ok = true
	}

	return
}

// TODO-P3 merge all the below
func (c *compiled) getSortedTypenExpanded(
	v string,
) (expandedActual []*sku.Transacted) {
	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)

	sa := iter.MakeFuncSetString(expandedMaybe)

	typExpander.Expand(sa, v)
	expandedActual = make([]*sku.Transacted, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			c.lock.Lock()
			defer c.lock.Unlock()

			ct, ok := c.Typen.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

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

func (c *compiled) getSortedEtikettenExpanded(
	v string,
) (expandedActual []*sku.Transacted) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)
	sa := iter.MakeFuncSetString(
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

func (c *compiled) getImplicitEtiketten(
	e *kennung.Etikett,
) kennung.EtikettSet {
	s, ok := c.ImplicitEtiketten[e.String()]

	if !ok || s == nil {
		return kennung.MakeEtikettSet()
	}

	return s
}

func (kc *Compiled) Cli() erworben.Cli {
	return kc.cli
}
