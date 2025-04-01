package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (repo *Repo) Checkin(
	skus sku.SkuTypeSetMutable,
	proto sku.Proto,
	delete bool,
	refreshCheckout bool,
) (processed sku.TransactedMutableSet, err error) {
	repo.Must(repo.Lock)

	processed = sku.MakeTransactedMutableSet()
	sortedResults := quiter.ElementsSorted(
		skus,
		func(left, right sku.SkuType) bool {
			return left.String() < right.String()
		},
	)

	for _, co := range sortedResults {
		if refreshCheckout {
			if err = repo.GetStore().GetStoreFS().RefreshCheckedOut(
				co,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		external := co.GetSkuExternal()

		if co.GetState() == checked_out_state.Untracked &&
			(co.GetSkuExternal().GetGenre() == genres.Zettel ||
				co.GetSkuExternal().GetGenre() == genres.Blob) {
			if external.Metadata.IsEmpty() {
				continue
			}

			if err = repo.GetStore().UpdateTransactedFromBlobs(
				co,
			); err != nil {
				if errors.Is(err, external_store.ErrUnsupportedOperation{}) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			external.ObjectId.Reset()

			if err = repo.GetStore().CreateOrUpdate(
				external,
				sku.StoreOptions{
					ApplyProto: true,
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if proto.Apply(external, genres.Zettel) {
				if err = repo.GetStore().CreateOrUpdate(
					external.GetSku(),
					sku.StoreOptions{},
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		} else {
			if err = repo.GetStore().CreateOrUpdateCheckedOut(
				co,
				!delete,
			); err != nil {
				err = errors.Wrapf(err, "CheckedOut: %s", co)
				return
			}
		}

		if !delete {
			continue
		}

		if err = repo.GetStore().DeleteCheckedOut(co); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = processed.Add(co.GetSkuExternal().CloneTransacted()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	repo.Must(repo.Unlock)

	return
}
