package env

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (u *Env) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		u.GetDirectoryLayout(),
		u.GetStore().GetBlobStore(),
		u.GetStore().GetStreamIndex(),
		u.MakeLuaVMPoolBuilder(),
		u,
	)
}

func (u *Env) MakeQueryBuilderExcludingHidden(
	dg ids.Genre,
) *query.Builder {
	if dg.IsEmpty() {
		dg = ids.MakeGenre(genres.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGenres(dg).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(u.GetConfig().GetFileExtensions()).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr()).
		WithHidden(u.GetMatcherArchiviert())
}

func (u *Env) MakeQueryBuilder(
	dg ids.Genre,
) *query.Builder {
	if dg.IsEmpty() {
		dg = ids.MakeGenre(genres.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGenres(dg).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(u.GetConfig().GetFileExtensions()).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr())
}

func (u *Env) GetDefaultExternalStore() *external_store.Store {
	return u.externalStores[ids.RepoId{}]
}
