package store_browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/external_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

type executor struct {
	store *Store
	qg    *query.Group
	out   interfaces.FuncIter[sku.SkuType]
	co    sku.CheckedOut
}

func (c *executor) tryToEmitOneExplicitlyCheckedOut(
	internal *sku.Transacted,
	item Item,
) (err error) {
	c.co.GetSkuExternal().ObjectId.Reset()

	var uSku *url.URL

	if uSku, err = c.store.getUrl(internal); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.TransactedResetter.ResetWith(c.co.GetSku(), internal)
	sku.TransactedResetter.ResetWith(c.co.GetSkuExternal().GetSku(), internal)

	if *uSku == item.Url.Url() {
		// c.co.SetState(checked_out_state.ExistsAndSame)
	} else {
		// c.co.SetState(checked_out_state.Changed)
	}

	c.co.GetSkuExternal().State = external_state.Tracked

	if err = c.tryToEmitOneCommon(item); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *executor) tryToEmitOneRecognized(
	internal *sku.Transacted,
	item Item,
) (err error) {
	c.co.SetState(checked_out_state.Recognized)

	if !c.qg.ContainsSkuCheckedOutState(c.co.GetState()) {
		return
	}

	sku.TransactedResetter.ResetWith(c.co.GetSku(), internal)
	sku.TransactedResetter.ResetWith(c.co.GetSkuExternal().GetSku(), internal)

	// if err = item.WriteToObjectId(&c.co.External.ObjectId); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	c.co.SetState(checked_out_state.Recognized)
	c.co.GetSkuExternal().State = external_state.Recognized

	if err = c.tryToEmitOneCommon(item); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *executor) tryToEmitOneUntracked(
	item Item,
) (err error) {
	c.co.SetState(checked_out_state.Untracked)

	if !c.qg.ContainsSkuCheckedOutState(c.co.GetState()) {
		return
	}

	sku.TransactedResetter.Reset(c.co.GetSkuExternal().GetSku())
	sku.TransactedResetter.Reset(c.co.GetSku())

	if err = c.co.GetSkuExternal().Metadata.Description.Set(item.Title); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.tryToEmitOneCommon(item); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *executor) tryToEmitOneCommon(
	i Item,
) (err error) {
	external := c.co.GetSkuExternal()

	if err = i.WriteToExternal(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	external.ObjectId.SetGenre(genres.Zettel)
	external.ExternalObjectId.SetGenre(genres.Zettel)

	if !c.qg.ContainsExternalSku(external, c.co.GetState()) {
		return
	}

	c.co.GetSkuExternal().RepoId = c.store.externalStoreInfo.RepoId

	if err = c.out(&c.co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
