package env

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (u *Env) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		u.GetFSHome(),
		u.GetStore().GetBlobStore(),
		u.GetStore().GetStreamIndex(),
		(&lua.VMPoolBuilder{}).WithSearcher(u.LuaSearcher),
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
		WithVirtualTags(u.config.Filters).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(u.GetConfig().FileExtensions).
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
		WithVirtualTags(u.config.Filters).
		WithRepoId(ids.RepoId{}).
		WithFileExtensionGetter(u.GetConfig().FileExtensions).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr())
}

func (u *Env) GetDefaultExternalStore() *external_store.Store {
	return u.externalStores[ids.RepoId{}]
}
