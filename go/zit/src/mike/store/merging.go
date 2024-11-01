package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) MergeCheckedOutIfNecessary(
	co *sku.CheckedOut,
) (commitOptions sku.CommitOptions, err error) {
	commitOptions.Mode = object_mode.ModeCommit

	if co.Internal.Metadata.Sha().IsNull() {
		return
	}

	var conflicts checkout_mode.Mode

	// TODO add checkout_mode.BlobOnly
	if co.Internal.Metadata.Sha().Equals(co.External.Metadata.Sha()) {
		commitOptions.Mode = object_mode.ModeEmpty
		return
	} else if co.Internal.Metadata.EqualsSansTai(&co.External.Metadata) {
		if !co.Internal.Metadata.Tai.Less(co.External.Metadata.Tai) {
      // TODO implement retroactive change
		}

    return
	} else if co.Internal.Metadata.Blob.Equals(&co.External.Metadata.Blob) {
		conflicts = checkout_mode.MetadataOnly
	} else {
		conflicts = checkout_mode.MetadataAndBlob
	}

	switch conflicts {
	case checkout_mode.BlobOnly:
	case checkout_mode.MetadataOnly:
	case checkout_mode.MetadataAndBlob:

	default:
	}

	co.State = checked_out_state.Conflicted

	return
}

func (s *Store) readExternalAndMergeIfNecessary(
	left, parent *sku.Transacted,
	options sku.CommitOptions,
) (err error) {
	if parent == nil {
		return
	}

	var col sku.CheckedOutLike

	if col, err = s.ReadCheckedOutFromTransacted(
		options.RepoId,
		parent,
	); err != nil {
		err = nil
		return
	}

	defer s.PutCheckedOutLike(col)

	right := col.GetSkuExternalLike().GetSku()

	parentEqualsExternal := right.Metadata.EqualsSansTai(&col.GetSku().Metadata)

	if parentEqualsExternal {
		op := checkout_options.OptionsWithoutMode{
			Force: true,
		}

		sku.TransactedResetter.ResetWithExceptFields(right, left)

		if err = s.UpdateCheckoutFromCheckedOut(
			op,
			col,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	tm := sku.Conflicted{
		CheckedOutLike: col,
		Left:           left,
		Middle:         parent,
		Right:          right,
	}

	if err = s.MergeConflicted(tm); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
