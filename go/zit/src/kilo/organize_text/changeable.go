package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func key(el sku.ExternalLike) string {
	eoid := el.GetExternalObjectId().String()
	if len(eoid) > 1 {
		return eoid
	}

	oid := el.GetSku().ObjectId.String()

	if len(oid) > 1 {
		return oid
	}

	desc := el.GetSku().Metadata.Description.String()
	if desc != "" {
		return desc
	}

	panic(fmt.Sprintf("empty key for external like: %#v", el))
}

// TODO explore using shas as keys
func keySha(el sku.ExternalLike) string {
	objectSha := &el.GetSku().Metadata.SelfMetadataWithoutTai

	if objectSha.IsNull() {
		panic("empty object sha")
	}

	return fmt.Sprintf(
		"%s.%s.%s",
		el.GetObjectId(),
		el.GetExternalObjectId(),
		objectSha,
	)
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

	for _, o := range a.All() {
		var selwi skuExternalLikeWithIndex
		var z sku.ExternalLike
		ok := false

		k := key(o.External)

		if selwi, ok = out.m[k]; !ok {
			z = organizeText.ObjectFactory.Get()

			organizeText.ObjectFactory.ResetWith(z, o.External)

			if !organizeText.Metadata.Type.IsEmpty() {
				z.GetSku().Metadata.Type.ResetWith(organizeText.Metadata.Type)
			}

			out.Add(z)

			zPrime, hasOriginal := original.Get(original.Key(o.External))

			if hasOriginal {
				z.GetSku().Metadata.Blob.ResetWith(&zPrime.GetSku().Metadata.Blob)
				z.GetSku().Metadata.Type.ResetWith(zPrime.GetSku().Metadata.Type)
			}

			m := z.GetSku().GetMetadata()

			for e := range organizeText.Metadata.AllPtr() {
        if o.Type == tag_paths.TypeUnknown {
          continue
        }

				if _, ok := m.Cache.TagPaths.All.ContainsComparer(
					catgut.ComparerString(e.String()),
				); ok {
					continue
				}

				z.GetSku().AddTagPtr(e)
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

		m := z.GetSku().GetMetadata()

		for e := range expanded.AllPtr() {
			m.AddTagPtr(e)
		}
	}

	for _, c := range a.Children {
		if err = c.addToSet(organizeText, out, original); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
