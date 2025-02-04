package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (repo *Repo) MakeExternalQueryGroup(
	metaBuilder query.BuilderOptions,
	externalQueryOptions sku.ExternalQueryOptions,
	args ...string,
) (queryGroup *query.Group, err error) {
	builder := repo.MakeQueryBuilderExcludingHidden(ids.MakeGenre(), metaBuilder)

	if queryGroup, err = builder.BuildQueryGroupWithRepoId(
		externalQueryOptions,
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	queryGroup.ExternalQueryOptions = externalQueryOptions

	return
}

func (repo *Repo) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		repo.GetEnvRepo(),
		repo.GetStore().GetTypedBlobStore(),
		repo.GetStore().GetStreamIndex(),
		repo.envLua.MakeLuaVMPoolBuilder(),
		repo,
	)
}

func (repo *Repo) MakeQueryBuilderExcludingHidden(
	genre ids.Genre,
	options query.BuilderOptions,
) *query.Builder {
	if genre.IsEmpty() {
		genre = ids.MakeGenre(genres.Zettel)
	}

	return repo.makeQueryBuilder().
		WithDefaultGenres(genre).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(repo.GetConfig().GetFileExtensions()).
		WithExpanders(repo.GetStore().GetAbbrStore().GetAbbr()).
		WithHidden(repo.GetMatcherDormant()).
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
