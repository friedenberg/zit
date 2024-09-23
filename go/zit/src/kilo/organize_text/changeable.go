package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func key(el sku.ExternalLike) string {
	eoid := el.GetExternalObjectId()
	if !eoid.IsEmpty() {
		return eoid.String()
	}

	if !el.GetSku().ObjectId.IsEmpty() {
		return el.GetSku().ObjectId.String()
	}

	desc := el.GetSku().Metadata.Description.String()
	if desc != "" {
		return desc
	}

	panic(fmt.Sprintf("empty key for external like: %#v", el))
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
	organizeText *Text,
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
			var selwi skuExternalLikeWithIndex
			var z sku.ExternalLike
			ok := false

			k := key(o.External)

			if selwi, ok = out.m[k]; !ok {
				z = organizeText.ObjectFactory.Get()

				organizeText.ObjectFactory.ResetWith(z, o.External)

				if err = organizeText.EachPtr(
					z.GetSku().AddTagPtr,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !organizeText.Metadata.Type.IsEmpty() {
					z.GetSku().Metadata.Type.ResetWith(organizeText.Metadata.Type)
				}

				out.Add(z)

				zPrime, hasOriginal := original.Get(original.Key(o.External))

				if hasOriginal {
					z.GetSku().Metadata.Blob.ResetWith(&zPrime.GetSku().Metadata.Blob)
					z.GetSku().Metadata.Type.ResetWith(zPrime.GetSku().Metadata.Type)
				}

				if !organizeText.Metadata.Type.IsEmpty() {
					z.GetSku().Metadata.Type.ResetWith(organizeText.Metadata.Type)
				}
			} else {
				z = selwi.ExternalLike
			}

			if o.External.GetSku().ObjectId.String() == "" {
				panic(fmt.Sprintf("%s: object id is nil", o))
			}

			if z == nil {
				panic("empty object")
			}

			if err = z.GetSku().Metadata.Description.Set(
				o.External.GetSku().Metadata.Description.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !o.External.GetSku().Metadata.Type.IsEmpty() {
				if err = z.GetSku().Metadata.Type.Set(
					o.External.GetSku().Metadata.Type.String(),
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
				o.External.GetSku().Metadata.Comments...,
			)

			if err = o.External.GetSku().Metadata.GetTags().EachPtr(
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
		if err = c.addToSet(organizeText, out, original); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
