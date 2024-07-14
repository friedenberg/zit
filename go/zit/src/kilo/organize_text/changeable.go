package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func key(sk *sku.Transacted) string {
	if sk.ObjectId.IsEmpty() {
		s := sk.Metadata.Description.String()

		if s == "" {
			panic("empty key")
		}

		return s
	} else {
		return sk.ObjectId.String()
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
					z.Metadata.Type.ResetWith(ot.Metadatei.Typ)
				}

				out.Add(z)

				zPrime, hasOriginal := original.Get(original.Key(&o.Transacted))

				if hasOriginal {
					z.Metadata.Blob.ResetWith(&zPrime.Metadata.Blob)
					z.Metadata.Type.ResetWith(zPrime.Metadata.Type)
				}

				if !ot.Metadatei.Typ.IsEmpty() {
					z.Metadata.Type.ResetWith(ot.Metadatei.Typ)
				}
			}

			if o.ObjectId.String() == "" {
				panic(fmt.Sprintf("%s: Kennung is nil", o))
			}

			if err = z.Metadata.Description.Set(
				o.Metadata.Description.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !o.Metadata.Type.IsEmpty() {
				if err = z.Metadata.Type.Set(
					o.Metadata.Type.String(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if !o.IsDirectOrSelf() {
				return
			}

			z.Metadata.Comments = append(
				z.Metadata.Comments,
				o.Metadata.Comments...,
			)

			if err = o.Metadata.GetTags().EachPtr(
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
