package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
)

func MakeBuilder(
	s fs_home.Home,
	blob_store *blob_store.VersionedStores,
	object_probe_index sku.ObjectProbeIndex,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	repoGetter sku.ExternalStoreForQueryGetter,
) (b *Builder) {
	b = &Builder{
		fs_home:                    s,
		blob_store:                 blob_store,
		object_probe_index:         object_probe_index,
		luaVMPoolBuilder:           luaVMPoolBuilder,
		virtualEtikettenBeforeInit: make(map[string]string),
		repoGetter:                 repoGetter,
	}

	return
}

type Builder struct {
	fs_home                    fs_home.Home
	blob_store                 *blob_store.VersionedStores
	object_probe_index         sku.ObjectProbeIndex
	luaVMPoolBuilder           *lua.VMPoolBuilder
	pinnedObjectIds            []ObjectId
	pinnedExternalObjectIds    []sku.ExternalObjectId
	repoGetter                 sku.ExternalStoreForQueryGetter
	repoId                     ids.RepoId
	virtualEtikettenBeforeInit map[string]string
	fileExtensionGetter        interfaces.FileExtensionGetter
	expanders                  ids.Abbr
	hidden                     sku.Query
	defaultGenres              ids.Genre
	defaultSigil               ids.Sigil
	permittedSigil             ids.Sigil
	doNotMatchEmpty            bool
	debug                      bool
	requireNonEmptyQuery       bool
}

func (b *Builder) makeState() *buildState {
	state := &buildState{
		builder:          b,
		latentErrors:     errors.MakeMulti(),
		virtualEtiketten: make(map[string]Lua),
	}

	if b.luaVMPoolBuilder != nil {
		state.luaVMPoolBuilder = b.luaVMPoolBuilder.Clone()
	}

	state.qg = state.makeGroup()

	state.pinnedObjectIds = make([]ObjectId, len(b.pinnedObjectIds))
	copy(state.pinnedObjectIds, b.pinnedObjectIds)

	state.pinnedExternalObjectIds = make(
		[]sku.ExternalObjectId,
		len(b.pinnedExternalObjectIds),
	)

	copy(state.pinnedExternalObjectIds, b.pinnedExternalObjectIds)

	return state
}

func (b *Builder) WithPermittedSigil(s ids.Sigil) *Builder {
	b.permittedSigil.Add(s)
	return b
}

func (b *Builder) WithDoNotMatchEmpty() *Builder {
	b.doNotMatchEmpty = true
	return b
}

func (b *Builder) WithRequireNonEmptyQuery() *Builder {
	b.requireNonEmptyQuery = true
	return b
}

func (mb *Builder) WithVirtualTags(vs map[string]string) *Builder {
	for k, v := range vs {
		mb.virtualEtikettenBeforeInit["%"+k] = v
	}

	return mb
}

func (mb *Builder) WithDebug() *Builder {
	mb.debug = true
	return mb
}

func (mb *Builder) WithRepoId(
	repoId ids.RepoId,
) *Builder {
	mb.repoId = repoId
	return mb
}

func (mb *Builder) WithFileExtensionGetter(
	feg interfaces.FileExtensionGetter,
) *Builder {
	mb.fileExtensionGetter = feg
	return mb
}

func (mb *Builder) WithExpanders(
	expanders ids.Abbr,
) *Builder {
	mb.expanders = expanders
	return mb
}

func (mb *Builder) WithDefaultGenres(
	defaultGenres ids.Genre,
) *Builder {
	mb.defaultGenres = defaultGenres
	return mb
}

func (mb *Builder) WithDefaultSigil(
	defaultSigil ids.Sigil,
) *Builder {
	mb.defaultSigil = defaultSigil
	return mb
}

func (mb *Builder) WithHidden(
	hidden sku.Query,
) *Builder {
	mb.hidden = hidden
	return mb
}

func (b *Builder) WithExternalLike(
	zts sku.ExternalLikeSet,
) *Builder {
	errors.PanicIfError(zts.Each(
		func(t sku.ExternalLike) (err error) {
			b.pinnedExternalObjectIds = append(
				b.pinnedExternalObjectIds,
				t.GetExternalObjectId(),
			)

			return
		},
	))

	return b
}

func (b *Builder) WithTransacted(
	zts sku.TransactedSet,
) *Builder {
	errors.PanicIfError(zts.Each(
		func(t *sku.Transacted) (err error) {
			b.pinnedObjectIds = append(
				b.pinnedObjectIds,
				ObjectId{
					ObjectIdLike: t.ObjectId.Clone(),
				},
			)

			return
		},
	))

	return b
}

func (b *Builder) WithCheckedOut(
	cos sku.CheckedOutLikeSet,
) *Builder {
	errors.PanicIfError(cos.Each(
		func(co sku.CheckedOutLike) (err error) {
			b.pinnedObjectIds = append(
				b.pinnedObjectIds,
				ObjectId{
					ObjectIdLike: co.GetSku().ObjectId.Clone(),
				},
			)

			return
		},
	))

	return b
}

func (b *Builder) WithOptionsFromOriginalQuery(
	qg *Group,
) *Builder {
	b.doNotMatchEmpty = !qg.matchOnEmpty
	return b
}

func (b *Builder) BuildQueryGroupWithRepoId(
	k ids.RepoId,
	eqo sku.ExternalQueryOptions,
	vs ...string,
) (qg *Group, err error) {
	state := b.makeState()

	ok := false
	state.eqo = eqo
	state.repo, ok = b.repoGetter.GetExternalStoreForQuery(k)

	state.qg.RepoId = k
	state.qg.ExternalQueryOptions = eqo

	if !ok {
		err = errors.Errorf("kasten not found: %q", k)
		return
	}

	if err = b.build(state, vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg = state.qg

	return
}

func (b *Builder) BuildQueryGroup(vs ...string) (qg *Group, err error) {
	state := b.makeState()

	if err = b.build(state, vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg = state.qg

	return
}

func (b *Builder) build(state *buildState, vs ...string) (err error) {
	var latent errors.Multi

	if err, latent = state.build(vs...); err != nil {
		if !errors.IsBadRequest(err) {
			latent.Add(errors.Wrapf(err, "Query String: %q", vs))
			err = latent
		}

		errors.Wrap(err)

		return
	}

	ui.Log().Print(state.qg.StringDebug())

	return
}
