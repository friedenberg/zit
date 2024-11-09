package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func key(el external_store.SkuType) string {
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
func keySha(el external_store.SkuType) string {
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
	original external_store.SkuTypeSet,
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
	output SkuMapWithOrder,
	objectsFromBefore external_store.SkuTypeSet,
) (err error) {
	expanded := ids.MakeTagMutableSet()

	if err = a.AllTags(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, organizeObject := range a.All() {
		var outputObject external_store.SkuType

		objectKey := key(organizeObject.sku)

		previouslyProcessedObject, wasPreviouslyProcessed := output.m[objectKey]

		if !wasPreviouslyProcessed {
			outputObject = ot.ObjectFactory.Get()

			ot.ObjectFactory.ResetWith(outputObject, organizeObject.sku)

			if !ot.Metadata.Type.IsEmpty() {
				outputObject.GetSku().Metadata.Type.ResetWith(ot.Metadata.Type)
			}

			outputObject.GetSku().RepoId.ResetWith(ot.Metadata.RepoId)

			output.Add(outputObject)

			objectOriginal, hasOriginal := objectsFromBefore.Get(
				objectsFromBefore.Key(organizeObject.sku),
			)

			if hasOriginal {
				outputObject.GetSku().Metadata.Blob.ResetWith(
					&objectOriginal.GetSku().Metadata.Blob,
				)

				outputObject.GetSku().Metadata.Type.ResetWith(
					objectOriginal.GetSku().Metadata.Type,
				)
			}

			outputMetadata := outputObject.GetSku().GetMetadata()

			for e := range ot.Metadata.AllPtr() {
				if organizeObject.tipe == tag_paths.TypeUnknown {
					continue
				}

				if _, ok := outputMetadata.Cache.TagPaths.All.ContainsComparer(
					catgut.ComparerString(e.String()),
				); ok {
					continue
				}

				outputObject.GetSku().AddTagPtr(e)
			}

			if !ot.Metadata.Type.IsEmpty() {
				outputObject.GetSku().Metadata.Type.ResetWith(ot.Metadata.Type)
			}
		} else {
			outputObject = previouslyProcessedObject.sku
		}

		if organizeObject.sku.GetSku().ObjectId.String() == "" {
			panic(fmt.Sprintf("%s: object id is nil", organizeObject))
		}

		if outputObject == nil {
			panic("empty object")
		}

		if err = outputObject.GetSku().Metadata.Description.Set(
			organizeObject.sku.GetSku().Metadata.Description.String(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !organizeObject.sku.GetSku().Metadata.Type.IsEmpty() {
			if err = outputObject.GetSku().Metadata.Type.Set(
				organizeObject.sku.GetSku().Metadata.Type.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !organizeObject.tipe.IsDirectOrSelf() {
			return
		}

		outputObject.GetSku().Metadata.Comments = append(
			outputObject.GetSku().Metadata.Comments,
			organizeObject.sku.GetSku().Metadata.Comments...,
		)

		if err = organizeObject.sku.GetSku().Metadata.GetTags().EachPtr(
			outputObject.GetSku().AddTagPtr,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		for e := range expanded.AllPtr() {
			outputObject.GetSku().AddTagPtr(e)
		}
	}

	for _, c := range a.Children {
		if err = c.addToSet(ot, output, objectsFromBefore); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
