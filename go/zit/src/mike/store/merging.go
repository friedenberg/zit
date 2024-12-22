package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO combine with other method in this file
func (s *Store) MergeCheckedOutIfNecessary(
	co *sku.CheckedOut,
	parentNegotiator sku.ParentNegotiator,
) (commitOptions sku.CommitOptions, err error) {
	commitOptions.Mode = object_mode.ModeCommit

	if co.GetSku().Metadata.Sha().IsNull() {
		return
	}

	var conflicts checkout_mode.Mode

	// TODO add checkout_mode.BlobOnly
	if co.GetSku().Metadata.Sha().Equals(co.GetSkuExternal().Metadata.Sha()) {
		commitOptions.Mode = object_mode.ModeEmpty
		return
	} else if co.GetSku().Metadata.EqualsSansTai(&co.GetSkuExternal().Metadata) {
		if !co.GetSku().Metadata.Tai.Less(co.GetSkuExternal().Metadata.Tai) {
			// TODO implement retroactive change
		}

		return
	} else if co.GetSku().Metadata.Blob.Equals(&co.GetSkuExternal().Metadata.Blob) {
		conflicts = checkout_mode.MetadataOnly
	} else {
		conflicts = checkout_mode.MetadataAndBlob
	}

	// TODO write conflicts
	switch conflicts {
	case checkout_mode.BlobOnly:
	case checkout_mode.MetadataOnly:
	case checkout_mode.MetadataAndBlob:
	default:
	}

	conflicted := sku.Conflicted{
		CheckedOut: co,
		Local:      co.GetSku(),
		Remote:     co.GetSkuExternal(),
	}

	if err = conflicted.FindBestCommonAncestor(parentNegotiator); err != nil {
		err = errors.Wrap(err)
		return
	}

	var skuReplacement *sku.Transacted

	// TODO pass mode / conflicts
	if skuReplacement, err = s.GetStoreFS().MakeMergedTransacted(
		conflicted,
	); err != nil {
		if sku.IsErrMergeConflict(err) {
			if err = s.GetStoreFS().GenerateConflictMarker(
				conflicted,
				conflicted.CheckedOut,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			co.SetState(checked_out_state.Conflicted)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	sku.TransactedResetter.ResetWith(co.GetSkuExternal(), skuReplacement)

	return
}

func (s *Store) ReadExternalAndMergeIfNecessary(
	left, parent *sku.Transacted,
	options sku.CommitOptions,
) (err error) {
	if parent == nil {
		return
	}

	var co *sku.CheckedOut

	if co, err = s.ReadCheckedOutFromTransacted(
		options.RepoId,
		parent,
	); err != nil {
		err = nil
		return
	}

	defer s.PutCheckedOutLike(co)

	right := co.GetSkuExternal().GetSku()

	parentEqualsExternal := right.Metadata.EqualsSansTai(&co.GetSku().Metadata)

	if parentEqualsExternal {
		op := checkout_options.OptionsWithoutMode{
			Force: true,
		}

		sku.TransactedResetter.ResetWithExceptFields(right, left)

		if err = s.UpdateCheckoutFromCheckedOut(
			op,
			co,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	conflicted := sku.Conflicted{
		CheckedOut: co,
		Local:      left,
		Base:       parent,
		Remote:     right,
	}

	if err = s.MergeConflicted(conflicted); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
