package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (repo *Repo) GetZettelFromObjectId(
	objectIdString string,
) (sk *sku.Transacted, err error) {
	builder := repo.MakeQueryBuilder(ids.MakeGenre(genres.Zettel), nil)

	var queryGroup *query.Query

	if queryGroup, err = builder.BuildQueryGroupWithRepoId(
		sku.ExternalQueryOptions{},
		objectIdString,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = repo.GetStore().QueryExactlyOneExternal(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (repo *Repo) GetObjectFromObjectId(
	objectIdString string,
) (sk *sku.Transacted, err error) {
	builder := repo.MakeQueryBuilder(ids.MakeGenre(genres.All()...), nil)

	var queryGroup *query.Query

	if queryGroup, err = builder.BuildQueryGroupWithRepoId(
		sku.ExternalQueryOptions{},
		objectIdString,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = repo.GetStore().QueryExactlyOneExternal(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
