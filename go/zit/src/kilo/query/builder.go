package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

func MakeBuilder(
	envRepo env_repo.Env,
	typedBlobStore typed_blob_store.Stores,
	objectProbeIndex sku.ObjectProbeIndex,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	repoGetter sku.ExternalStoreForQueryGetter,
) (b *Builder) {
	b = &Builder{
		envRepo:          envRepo,
		typedBlobStore:   typedBlobStore,
		objectProbeIndex: objectProbeIndex,
		luaVMPoolBuilder: luaVMPoolBuilder,
		repoGetter:       repoGetter,
	}

	return
}

type Builder struct {
	envRepo                 env_repo.Env
	typedBlobStore          typed_blob_store.Stores
	objectProbeIndex        sku.ObjectProbeIndex
	luaVMPoolBuilder        *lua.VMPoolBuilder
	pinnedObjectIds         []pinnedObjectId
	pinnedExternalObjectIds []sku.ExternalObjectId
	repoGetter              sku.ExternalStoreForQueryGetter
	repoId                  ids.RepoId
	fileExtensionGetter     interfaces.FileExtensionGetter
	expanders               ids.Abbr
	hidden                  sku.Query
	doNotMatchEmpty         bool
	debug                   bool
	requireNonEmptyQuery    bool
	defaultQuery            string
	workspaceEnabled        bool

	options options
}

func (builder *Builder) makeState() *buildState {
	state := &buildState{
		options:      builder.options,
		builder:      builder,
		latentErrors: errors.MakeMulti(),
	}

	if builder.luaVMPoolBuilder != nil {
		state.luaVMPoolBuilder = builder.luaVMPoolBuilder.Clone()
	}

	state.group = state.makeGroup()

	state.pinnedObjectIds = make([]pinnedObjectId, len(builder.pinnedObjectIds))
	copy(state.pinnedObjectIds, builder.pinnedObjectIds)

	state.pinnedExternalObjectIds = make(
		[]sku.ExternalObjectId,
		len(builder.pinnedExternalObjectIds),
	)

	copy(state.pinnedExternalObjectIds, builder.pinnedExternalObjectIds)

	return state
}

func (b *Builder) WithOptions(options BuilderOption) *Builder {
	if options == nil {
		return b
	}

	applied := options.Apply(b)

	if applied != nil {
		return applied
	}

	return b
}

// TODO refactor into BuilderOption
func (b *Builder) WithPermittedSigil(s ids.Sigil) *Builder {
	b.options.permittedSigil.Add(s)
	return b
}

// TODO refactor into BuilderOption
func (b *Builder) WithDoNotMatchEmpty() *Builder {
	b.doNotMatchEmpty = true
	return b
}

// TODO refactor into BuilderOption
func (b *Builder) WithRequireNonEmptyQuery() *Builder {
	b.requireNonEmptyQuery = true
	return b
}

// TODO refactor into BuilderOption
func (mb *Builder) WithDebug() *Builder {
	mb.debug = true
	return mb
}

// TODO refactor into BuilderOption
func (mb *Builder) WithRepoId(
	repoId ids.RepoId,
) *Builder {
	mb.repoId = repoId
	return mb
}

// TODO refactor into BuilderOption
func (mb *Builder) WithFileExtensionGetter(
	feg interfaces.FileExtensionGetter,
) *Builder {
	mb.fileExtensionGetter = feg
	return mb
}

// TODO refactor into BuilderOption
func (mb *Builder) WithExpanders(
	expanders ids.Abbr,
) *Builder {
	mb.expanders = expanders
	return mb
}

// TODO refactor into BuilderOption
func (mb *Builder) WithDefaultGenres(
	defaultGenres ids.Genre,
) *Builder {
	mb.options.defaultGenres = defaultGenres
	return mb
}

// TODO refactor into BuilderOption
func (mb *Builder) WithDefaultSigil(
	defaultSigil ids.Sigil,
) *Builder {
	mb.options.defaultSigil = defaultSigil
	return mb
}

func (mb *Builder) WithHidden(
	hidden sku.Query,
) *Builder {
	mb.hidden = hidden
	return mb
}

// TODO
func (b *Builder) WithExternalLike(
	zts sku.SkuTypeSet,
) *Builder {
	for t := range zts.All() {
		if t.GetExternalObjectId().IsEmpty() {
			b.pinnedObjectIds = append(
				b.pinnedObjectIds,
				pinnedObjectId{
					Sigil: ids.SigilExternal,
					ObjectId: ObjectId{
						Exact:    true,
						ObjectId: t.GetObjectId(),
					},
				},
			)
		} else {
			if t.GetExternalObjectId().GetGenre() == genres.None {
				panic(
					errors.BadRequestf(
						"External object ID has an empty genre: %q",
						t.GetExternalObjectId(),
					),
				)
			}

			b.pinnedExternalObjectIds = append(
				b.pinnedExternalObjectIds,
				t.GetExternalObjectId(),
			)
		}
	}

	return b
}

func (b *Builder) WithTransacted(
	zts sku.TransactedSet,
	sigil ids.Sigil,
) *Builder {
	errors.PanicIfError(zts.Each(
		func(t *sku.Transacted) (err error) {
			b.pinnedObjectIds = append(
				b.pinnedObjectIds,
				pinnedObjectId{
					Sigil: sigil,
					ObjectId: ObjectId{
						ObjectId: t.ObjectId.Clone(),
					},
				},
			)

			return
		},
	))

	return b
}

func (b *Builder) BuildQueryGroupWithRepoId(
	externalQueryOptions sku.ExternalQueryOptions,
	values ...string,
) (group *Query, err error) {
	state := b.makeState()

	ok := false
	state.repo, ok = b.repoGetter.GetExternalStoreForQuery(
		externalQueryOptions.RepoId,
	)

	state.group.RepoId = externalQueryOptions.RepoId
	state.group.ExternalQueryOptions = externalQueryOptions

	if !ok {
		err = errors.Errorf("kasten not found: %q", externalQueryOptions.RepoId)
		return
	}

	if err = b.build(state, values...); err != nil {
		err = errors.Wrap(err)
		return
	}

	group = state.group

	return
}

func (b *Builder) BuildQueryGroup(vs ...string) (group *Query, err error) {
	state := b.makeState()

	if err = b.build(state, vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	group = state.group

	return
}

func (b *Builder) build(state *buildState, values ...string) (err error) {
	var latent errors.Multi

	if err, latent = state.build(values...); err != nil {
		if !errors.IsBadRequest(err) {
			latent.Add(errors.Wrapf(err, "Query String: %q", values))
			err = latent
		}

		errors.Wrap(err)

		return
	}

	if len(state.missingBlobs) > 0 {
		me := errors.MakeMulti()

		for _, e := range state.missingBlobs {
			me.Add(e)
		}

		err = me

		return
	}

	if b.defaultQuery == "" {
		return
	}

	defaultQueryGroupState := state.copy()
	defaultQueryGroupState.options.defaultGenres = ids.MakeGenre(genres.All()...)

	if err, _ = defaultQueryGroupState.build(b.defaultQuery); err != nil {
		err = errors.Wrap(err)
		return
	}

	state.group.defaultQuery = defaultQueryGroupState.group

	ui.Log().Print(state.group.StringDebug())

	return
}
