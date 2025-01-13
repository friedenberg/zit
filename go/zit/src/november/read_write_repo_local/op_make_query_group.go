package read_write_repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (u *Repo) MakeExternalQueryGroup(
	metaBuilder query.BuilderOptions,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (qg *query.Group, err error) {
	b := u.MakeQueryBuilderExcludingHidden(ids.MakeGenre(), metaBuilder)

	if qg, err = b.BuildQueryGroupWithRepoId(
		externalQueryOptions,
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.ExternalQueryOptions = externalQueryOptions

	return
}

func (u *Repo) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		u.GetRepoLayout(),
		u.GetStore().GetBlobStore(),
		u.GetStore().GetStreamIndex(),
		u.MakeLuaVMPoolBuilder(),
		u,
	)
}

func (u *Repo) MakeQueryBuilderExcludingHidden(
	dg ids.Genre,
	options query.BuilderOptions,
) *query.Builder {
	if dg.IsEmpty() {
		dg = ids.MakeGenre(genres.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGenres(dg).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(u.GetConfig().GetFileExtensions()).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr()).
		WithHidden(u.GetMatcherDormant()).
		WithOptions(options)
}

func (u *Repo) MakeQueryBuilder(
	dg ids.Genre,
	options query.BuilderOptions,
) *query.Builder {
	if dg.IsEmpty() {
		dg = ids.MakeGenre(genres.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGenres(dg).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(u.GetConfig().GetFileExtensions()).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr()).
		WithOptions(options)
}

func (u *Repo) GetDefaultExternalStore() *external_store.Store {
	return u.externalStores[ids.RepoId{}]
}
