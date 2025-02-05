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
	defaultGenres           ids.Genre
	defaultSigil            ids.Sigil
	permittedSigil          ids.Sigil
	doNotMatchEmpty         bool
	debug                   bool
	requireNonEmptyQuery    bool
	defaultQuery            string
}

func (builder *Builder) makeState() *buildState {
	state := &buildState{
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

func (b *Builder) WithOptions(options BuilderOptions) *Builder {
	if options != nil {
		b = options.Apply(b)
	}

	return b
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

func (b *Builder) WithCheckedOut(
	cos sku.SkuTypeSet,
) *Builder {
	for co := range cos.All() {
		b.pinnedObjectIds = append(
			b.pinnedObjectIds,
			pinnedObjectId{
				Sigil: ids.SigilExternal,
				ObjectId: ObjectId{
					Exact:    true,
					ObjectId: co.GetSku().ObjectId.Clone(),
				},
			},
		)
	}

	return b
}

func (b *Builder) WithOptionsFromOriginalQuery(
	qg *Group,
) *Builder {
	b.doNotMatchEmpty = !qg.matchOnEmpty
	return b
}

func (b *Builder) BuildQueryGroupWithRepoId(
	externalQueryOptions sku.ExternalQueryOptions,
	vs ...string,
) (group *Group, err error) {
	state := b.makeState()

	ok := false
	state.repo, ok = b.repoGetter.GetExternalStoreForQuery(externalQueryOptions.RepoId)

	state.group.RepoId = externalQueryOptions.RepoId
	state.group.ExternalQueryOptions = externalQueryOptions

	if !ok {
		err = errors.Errorf("kasten not found: %q", externalQueryOptions.RepoId)
		return
	}

	if err = b.build(state, vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	group = state.group

	return
}

func (b *Builder) BuildQueryGroup(vs ...string) (group *Group, err error) {
	state := b.makeState()

	if err = b.build(state, vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	group = state.group

	return
}

func (b *Builder) build(state *buildState, args ...string) (err error) {
	if b.defaultQuery != "" {
		args = append(args, b.defaultQuery)
		// defaultQueryGroupState := state.copy()

		// if err, _ = defaultQueryGroupState.build(b.defaultQuery); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }
	}

	var latent errors.Multi

	if err, latent = state.build(args...); err != nil {
		if !errors.IsBadRequest(err) {
			latent.Add(errors.Wrapf(err, "Query String: %q", args))
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

	ui.Log().Print(state.group.StringDebug())

	return
}
