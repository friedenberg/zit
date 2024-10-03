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
	ot *Text,
	output SkuMapWithOrder,
	objectsFromBefore sku.ExternalLikeSet,
) (err error) {
	expanded := ids.MakeTagMutableSet()

	if err = a.AllTags(expanded); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, organizeObject := range a.All() {
		var outputObject sku.ExternalLike

		objectKey := key(organizeObject.External)

		previouslyProcessedObject, wasPreviouslyProcessed := output.m[objectKey]

		if !wasPreviouslyProcessed {
			outputObject = ot.ObjectFactory.Get()

			ot.ObjectFactory.ResetWith(outputObject, organizeObject.External)

			if !ot.Metadata.Type.IsEmpty() {
				outputObject.GetSku().Metadata.Type.ResetWith(ot.Metadata.Type)
			}

			outputObject.GetSku().RepoId.ResetWith(ot.Metadata.RepoId)

			output.Add(outputObject)

			objectOriginal, hasOriginal := objectsFromBefore.Get(
				objectsFromBefore.Key(organizeObject.External),
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
				if organizeObject.Type == tag_paths.TypeUnknown {
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
			outputObject = previouslyProcessedObject.ExternalLike
		}

		if organizeObject.External.GetSku().ObjectId.String() == "" {
			panic(fmt.Sprintf("%s: object id is nil", organizeObject))
		}

		if outputObject == nil {
			panic("empty object")
		}

		if err = outputObject.GetSku().Metadata.Description.Set(
			organizeObject.External.GetSku().Metadata.Description.String(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !organizeObject.External.GetSku().Metadata.Type.IsEmpty() {
			if err = outputObject.GetSku().Metadata.Type.Set(
				organizeObject.External.GetSku().Metadata.Type.String(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !organizeObject.IsDirectOrSelf() {
			return
		}

		outputObject.GetSku().Metadata.Comments = append(
			outputObject.GetSku().Metadata.Comments,
			organizeObject.External.GetSku().Metadata.Comments...,
		)

		if err = organizeObject.External.GetSku().Metadata.GetTags().EachPtr(
			outputObject.GetSku().AddTagPtr,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		m := outputObject.GetSku().GetMetadata()

		for e := range expanded.AllPtr() {
			m.AddTagPtr(e)
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
