package store_browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

type executor struct {
	store *Store
	qg    *query.Group
	out   interfaces.FuncIter[sku.CheckedOutLike]
	co    sku.CheckedOut
}

func (c *executor) tryToEmitOneExplicitlyCheckedOut(
	internal *sku.Transacted,
	item Item,
) (err error) {
	c.co.External.ObjectId.Reset()

	var uSku *url.URL

	if uSku, err = c.store.getUrl(internal); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.TransactedResetter.ResetWith(&c.co.Internal, internal)
	sku.TransactedResetter.ResetWith(c.co.External.GetSku(), internal)

	if *uSku == item.Url.URL {
		c.co.State = checked_out_state.ExistsAndSame
	} else {
		c.co.State = checked_out_state.ExistsAndDifferent
	}

	c.co.External.State = external_state.Tracked

	if err = c.tryToEmitOneCommon(item, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *executor) tryToEmitOneRecognized(
	internal *sku.Transacted,
	item Item,
) (err error) {
	c.co.State = checked_out_state.Recognized

	if !c.qg.ContainsSkuCheckedOutState(c.co.State) {
		return
	}

	sku.TransactedResetter.ResetWith(&c.co.Internal, internal)
	sku.TransactedResetter.ResetWith(c.co.External.GetSku(), internal)

	// if err = item.WriteToObjectId(&c.co.External.ObjectId); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	c.co.State = checked_out_state.Recognized
	c.co.External.State = external_state.Recognized

	if err = c.tryToEmitOneCommon(item, true); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *executor) tryToEmitOneUntracked(
	item Item,
) (err error) {
	c.co.State = checked_out_state.Untracked

	if !c.qg.ContainsSkuCheckedOutState(c.co.State) {
		return
	}

	sku.TransactedResetter.Reset(c.co.External.GetSku())
	sku.TransactedResetter.Reset(&c.co.Internal)

	if err = c.co.External.Metadata.Description.Set(item.Title); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.co.External.State = external_state.Untracked

	if err = c.tryToEmitOneCommon(item, true); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *executor) tryToEmitOneCommon(
	i Item,
	overwrite bool,
) (err error) {
	external := &c.co.External

	if err = i.WriteToExternal(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	external.ObjectId.SetGenre(genres.Zettel)

	if !c.qg.ContainsExternalSku(external, c.co.State) {
		return
	}

	if err = c.co.External.ObjectId.SetRepoId(
		c.store.externalStoreInfo.RepoId.String(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.out(&c.co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
