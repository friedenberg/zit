package local_working_copy

import (
	"fmt"

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

	workspaceConfig := repo.envWorkspace.GetWorkspaceConfig()

	if workspaceConfig != nil {
		defaultQueryGroup := workspaceConfig.GetDefaultQueryGroup()

    // TODO add after parsing as an independent query group, rather than as a
    // literal
		if defaultQueryGroup != "" {
			args = append(
				args,
				fmt.Sprintf("[%s]", workspaceConfig.GetDefaultQueryGroup()),
			)
		}
	}

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

func (u *Repo) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		u.GetEnvRepo(),
		u.GetStore().GetTypedBlobStore(),
		u.GetStore().GetStreamIndex(),
		u.envLua.MakeLuaVMPoolBuilder(),
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
