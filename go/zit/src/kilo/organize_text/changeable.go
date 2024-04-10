package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func key(sk *sku.Transacted) string {
	if sk.Kennung.IsEmpty() {
		s := sk.Metadatei.Bezeichnung.String()

		if s == "" {
			panic("empty key")
		}

		return s
	} else {
		return sk.Kennung.String()
	}
}

func (ot *Text) GetSkus(
	original sku.TransactedSet,
) (out map[string]*sku.Transacted, err error) {
	out = make(map[string]*sku.Transacted)

	if err = ot.addToSet(
		ot,
		out,
		original,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Assignment) addToSet(
	ot *Text,
	out map[string]*sku.Transacted,
	original sku.TransactedSet,
) (err error) {
	expanded := kennung.MakeEtikettMutableSet()

	if err = a.AllEtiketten(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Named.Each(
		func(o *obj) (err error) {
			var z *sku.Transacted
			ok := false

			if z, ok = out[key(&o.Transacted)]; !ok {
				z = sku.GetTransactedPool().Get()

				if err = z.SetFromSkuLike(&o.Transacted); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = ot.EachPtr(
					z.Metadatei.AddEtikettPtr,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadatei.Typ.ResetWith(ot.Metadatei.Typ)
				}

				out[key(z)] = z

				zPrime, hasOriginal := original.Get(original.Key(&o.Transacted))

				if hasOriginal {
					z.Metadatei.Akte.ResetWith(&zPrime.Metadatei.Akte)
					z.Metadatei.Typ.ResetWith(zPrime.Metadatei.Typ)
				}

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadatei.Typ.ResetWith(ot.Metadatei.Typ)
				}
			}

			if o.Kennung.String() == "" {
				panic(fmt.Sprintf("%s: Kennung is nil", o))
			}

			if err = z.Metadatei.Bezeichnung.Set(
				o.Metadatei.Bezeichnung.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !o.Metadatei.Typ.IsEmpty() {
				if err = z.Metadatei.Typ.Set(
					o.Metadatei.Typ.String(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			z.Metadatei.Comments = append(
				z.Metadatei.Comments,
				o.Metadatei.Comments...,
			)

			if err = o.Metadatei.GetEtiketten().EachPtr(
				z.Metadatei.AddEtikettPtr,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = expanded.EachPtr(
				z.Metadatei.AddEtikettPtr,
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

	if err = a.Unnamed.Each(
		func(o *obj) (err error) {
			var z *sku.Transacted
			ok := false

			if z, ok = out[key(&o.Transacted)]; !ok {
				z = sku.GetTransactedPool().Get()

				if err = z.SetFromSkuLike(&o.Transacted); err != nil {
					err = errors.Wrap(err)
					return
				}

				z.Kennung.SetGattung(gattung.Zettel)

				if err = ot.EachPtr(
					z.Metadatei.AddEtikettPtr,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadatei.Typ.ResetWith(ot.Metadatei.Typ)
				}

				out[key(z)] = z

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadatei.Typ.ResetWith(ot.Metadatei.Typ)
				}
			}

			if err = z.Metadatei.Bezeichnung.Set(
				o.Metadatei.Bezeichnung.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !o.Metadatei.Typ.IsEmpty() {
				if err = z.Metadatei.Typ.Set(
					o.Metadatei.Typ.String(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			z.Metadatei.Comments = append(
				z.Metadatei.Comments,
				o.Metadatei.Comments...,
			)

			if err = o.Metadatei.GetEtiketten().EachPtr(
				z.Metadatei.AddEtikettPtr,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = expanded.EachPtr(
				z.Metadatei.AddEtikettPtr,
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

	for _, c := range a.Children {
		if err = c.addToSet(ot, out, original); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
