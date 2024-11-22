package organize_text

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

var keyer = sku.GetExternalLikeKeyer[sku.SkuType]()

func key(el sku.SkuType) string {
	key := keyer.GetKey(el)
	return key
}

// TODO explore using shas as keys
func keySha(el sku.SkuType) string {
	objectSha := &el.GetSkuExternal().Metadata.SelfMetadataWithoutTai

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
	original sku.SkuTypeSet,
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
	objectsFromBefore sku.SkuTypeSet,
) (err error) {
	expanded := ids.MakeTagMutableSet()

	if err = a.AllTags(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, organizeObject := range a.All() {
		var outputObject sku.SkuType

		objectKey := key(organizeObject.sku)

		previouslyProcessedObject, wasPreviouslyProcessed := output.m[objectKey]

		if !wasPreviouslyProcessed {
			outputObject = ot.ObjectFactory.Get()

			ot.ObjectFactory.ResetWith(outputObject, organizeObject.sku)

			if !ot.Metadata.Type.IsEmpty() {
				outputObject.GetSkuExternal().Metadata.Type.ResetWith(ot.Metadata.Type)
			}

			outputObject.GetSkuExternal().RepoId.ResetWith(ot.Metadata.RepoId)

			output.Add(outputObject)

			objectOriginal, hasOriginal := objectsFromBefore.Get(objectKey)

			if hasOriginal {
				outputObject.GetSkuExternal().Metadata.Blob.ResetWith(
					&objectOriginal.GetSkuExternal().Metadata.Blob,
				)

				outputObject.GetSkuExternal().Metadata.Type.ResetWith(
					objectOriginal.GetSkuExternal().Metadata.Type,
				)

				outputObject.GetSkuExternal().GetSkuExternal().Metadata.Blob.ResetWith(
					&objectOriginal.GetSkuExternal().GetSkuExternal().Metadata.Blob,
				)

				outputObject.GetSkuExternal().GetSkuExternal().Metadata.Type.ResetWith(
					objectOriginal.GetSkuExternal().GetSkuExternal().Metadata.Type,
				)

				outputObject.SetState(objectOriginal.GetState())

				{
					src := &objectOriginal.GetSkuExternal().Metadata
					dst := &outputObject.GetSkuExternal().Metadata
					dst.Fields = objectOriginal.GetSkuExternal().Metadata.Fields[:0]
					dst.Fields = append(dst.Fields, src.Fields...)
				}
			}

			outputMetadata := outputObject.GetSkuExternal().GetMetadata()

			for e := range ot.Metadata.AllPtr() {
				if organizeObject.tipe == tag_paths.TypeUnknown {
					continue
				}

				if _, ok := outputMetadata.Cache.TagPaths.All.ContainsComparer(
					catgut.ComparerString(e.String()),
				); ok {
					continue
				}

				outputObject.GetSkuExternal().AddTagPtr(e)
			}

			if !ot.Metadata.Type.IsEmpty() {
				outputObject.GetSkuExternal().Metadata.Type.ResetWith(ot.Metadata.Type)
			}
		} else {
			outputObject = previouslyProcessedObject.sku
		}

		if organizeObject.GetSkuExternal().ObjectId.String() == "" {
			panic(fmt.Sprintf("%s: object id is nil", organizeObject))
		}

		if outputObject == nil {
			panic("empty object")
		}

		if err = outputObject.GetSkuExternal().Metadata.Description.Set(
			organizeObject.GetSkuExternal().Metadata.Description.String(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !organizeObject.GetSkuExternal().Metadata.Type.IsEmpty() {
			if err = outputObject.GetSkuExternal().Metadata.Type.Set(
				organizeObject.GetSkuExternal().Metadata.Type.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !organizeObject.tipe.IsDirectOrSelf() {
			return
		}

		outputObject.GetSkuExternal().Metadata.Comments = append(
			outputObject.GetSkuExternal().Metadata.Comments,
			organizeObject.GetSkuExternal().Metadata.Comments...,
		)

		if err = organizeObject.GetSkuExternal().Metadata.GetTags().EachPtr(
			outputObject.GetSkuExternal().AddTagPtr,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		for e := range expanded.AllPtr() {
			outputObject.GetSkuExternal().AddTagPtr(e)
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
