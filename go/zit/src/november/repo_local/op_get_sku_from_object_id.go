package repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (u *Repo) GetSkuFromObjectId(
	objectIdString string,
) (sk *sku.Transacted, err error) {
	b := u.MakeQueryBuilder(ids.MakeGenre(genres.Zettel))

	var qg *query.Group

	if qg, err = b.BuildQueryGroup(objectIdString); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = u.GetStore().QueryExactlyOne(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
