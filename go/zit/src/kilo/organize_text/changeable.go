package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func key(sk skuType) string {
	if sk.GetSku().ObjectId.IsEmpty() {
		s := sk.GetSku().Metadata.Description.String()

		if s == "" {
			panic("empty key")
		}

		return s
	} else {
		return sk.GetSku().ObjectId.String()
	}
}

func (ot *Text) GetSkus(
	original sku.ExternalLikeSet,
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
	original sku.ExternalLikeSet,
) (err error) {
	expanded := ids.MakeTagMutableSet()

	if err = a.AllTags(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.Each(
		func(o *obj) (err error) {
			var z skuType
			ok := false

			if z, ok = out.m[key(o.Transacted)]; !ok {
				z = sku.GetTransactedPool().Get()

				if err = z.GetSku().SetFromSkuLike(o.Transacted.GetSku()); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = ot.EachPtr(
					z.GetSku().AddTagPtr,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !ot.Metadata.Typ.IsEmpty() {
					z.GetSku().Metadata.Type.ResetWith(ot.Metadata.Typ)
				}

				out.Add(z)

				zPrime, hasOriginal := original.Get(original.Key(o.Transacted))

				if hasOriginal {
					z.GetSku().Metadata.Blob.ResetWith(&zPrime.GetSku().Metadata.Blob)
					z.GetSku().Metadata.Type.ResetWith(zPrime.GetSku().Metadata.Type)
				}

				if !ot.Metadata.Typ.IsEmpty() {
					z.GetSku().Metadata.Type.ResetWith(ot.Metadata.Typ)
				}
			}

			if o.Transacted.GetSku().ObjectId.String() == "" {
				panic(fmt.Sprintf("%s: object id is nil", o))
			}

			if err = z.GetSku().Metadata.Description.Set(
				o.Transacted.GetSku().Metadata.Description.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !o.Transacted.GetSku().Metadata.Type.IsEmpty() {
				if err = z.GetSku().Metadata.Type.Set(
					o.Transacted.GetSku().Metadata.Type.String(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if !o.IsDirectOrSelf() {
				return
			}

			z.GetSku().Metadata.Comments = append(
				z.GetSku().Metadata.Comments,
				o.Transacted.GetSku().Metadata.Comments...,
			)

			if err = o.Transacted.GetSku().Metadata.GetTags().EachPtr(
				z.GetSku().AddTagPtr,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = expanded.EachPtr(
				z.GetSku().AddTagPtr,
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
