package umwelt

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/external_store"
)

func (u *Umwelt) makeQueryBuilder() *query.Builder {
	return query.MakeBuilder(
		u.Standort(),
		u.GetStore().GetAkten(),
		u.GetStore().GetVerzeichnisse(),
		(&lua.VMPoolBuilder{}).WithSearcher(u.LuaSearcher),
		u,
	)
}

func (u *Umwelt) MakeQueryBuilderExcludingHidden(
	dg kennung.Genre,
) *query.Builder {
	if dg.IsEmpty() {
		dg = kennung.MakeGenre(gattung.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGattungen(dg).
		WithVirtualEtiketten(u.konfig.Filters).
		WithKasten(u.GetDefaultExternalStore()).
		WithFileExtensionGetter(u.GetKonfig().FileExtensions).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr()).
		WithHidden(u.GetMatcherArchiviert())
}

func (u *Umwelt) MakeQueryBuilder(
	dg kennung.Genre,
) *query.Builder {
	if dg.IsEmpty() {
		dg = kennung.MakeGenre(gattung.Zettel)
	}

	return u.makeQueryBuilder().
		WithDefaultGattungen(dg).
		WithVirtualEtiketten(u.konfig.Filters).
		WithKasten(u.GetDefaultExternalStore()).
		WithFileExtensionGetter(u.GetKonfig().FileExtensions).
		WithExpanders(u.GetStore().GetAbbrStore().GetAbbr())
}

func (u *Umwelt) GetDefaultExternalStore() *external_store.Store {
	return u.externalStores[""]
}