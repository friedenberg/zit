package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func key(sk *sku.Transacted) string {
	if sk.Kennung.IsEmpty() {
		s := sk.Metadatei.Description.String()

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
) (out SkuMapWithOrder, err error) {
	out = MakeSkuMapWithOrder(original.Len())

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
	out SkuMapWithOrder,
	original sku.TransactedSet,
) (err error) {
	expanded := ids.MakeTagMutableSet()

	if err = a.AllEtiketten(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Each(
		func(o *obj) (err error) {
			var z *sku.Transacted
			ok := false

			if z, ok = out.m[key(&o.Transacted)]; !ok {
				z = sku.GetTransactedPool().Get()

				if err = z.SetFromSkuLike(&o.Transacted); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = ot.EachPtr(
					z.AddTagPtr,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadatei.Type.ResetWith(ot.Metadatei.Typ)
				}

				out.Add(z)

				zPrime, hasOriginal := original.Get(original.Key(&o.Transacted))

				if hasOriginal {
					z.Metadatei.Blob.ResetWith(&zPrime.Metadatei.Blob)
					z.Metadatei.Type.ResetWith(zPrime.Metadatei.Type)
				}

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadatei.Type.ResetWith(ot.Metadatei.Typ)
				}
			}

			if o.Kennung.String() == "" {
				panic(fmt.Sprintf("%s: Kennung is nil", o))
			}

			if err = z.Metadatei.Description.Set(
				o.Metadatei.Description.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !o.Metadatei.Type.IsEmpty() {
				if err = z.Metadatei.Type.Set(
					o.Metadatei.Type.String(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if !o.IsDirectOrSelf() {
				return
			}

			z.Metadatei.Comments = append(
				z.Metadatei.Comments,
				o.Metadatei.Comments...,
			)

			if err = o.Metadatei.GetTags().EachPtr(
				z.AddTagPtr,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = expanded.EachPtr(
				z.AddTagPtr,
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
