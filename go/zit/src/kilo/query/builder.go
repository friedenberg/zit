package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/india/env_workspace"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

func MakeBuilder(
	s env_repo.Env,
	typedBlobStore typed_blob_store.Stores,
	objectProbeIndex sku.ObjectProbeIndex,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	repoGetter sku.ExternalStoreForQueryGetter,
	envWorkspace env_workspace.Env,
) (b *Builder) {
	b = &Builder{
		envRepo:          s,
		typedBlobStore:   typedBlobStore,
		objectProbeIndex: objectProbeIndex,
		luaVMPoolBuilder: luaVMPoolBuilder,
		repoGetter:       repoGetter,
		envWorkspace:     envWorkspace,
	}

	return
}

type Builder struct {
	envRepo                 env_repo.Env
	envWorkspace            env_workspace.Env
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
}

func (b *Builder) makeState() *buildState {
	state := &buildState{
		builder:      b,
		latentErrors: errors.MakeMulti(),
	}

	if b.luaVMPoolBuilder != nil {
		state.luaVMPoolBuilder = b.luaVMPoolBuilder.Clone()
	}

	state.qg = state.makeGroup()

	state.pinnedObjectIds = make([]pinnedObjectId, len(b.pinnedObjectIds))
	copy(state.pinnedObjectIds, b.pinnedObjectIds)

	state.pinnedExternalObjectIds = make(
		[]sku.ExternalObjectId,
		len(b.pinnedExternalObjectIds),
	)

	copy(state.pinnedExternalObjectIds, b.pinnedExternalObjectIds)

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
	eqo sku.ExternalQueryOptions,
	vs ...string,
) (qg *Group, err error) {
	state := b.makeState()

	ok := false
	state.eqo = eqo
	state.repo, ok = b.repoGetter.GetExternalStoreForQuery(eqo.RepoId)

	state.qg.RepoId = eqo.RepoId
	state.qg.ExternalQueryOptions = eqo

	if !ok {
		err = errors.Errorf("kasten not found: %q", eqo.RepoId)
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

func (b *Builder) build(state *buildState, args ...string) (err error) {
	if b.envWorkspace != nil {
		workspaceConfig := b.envWorkspace.GetWorkspaceConfig()

		if workspaceConfig != nil {
			defaultQueryGroup := workspaceConfig.GetDefaultQueryGroup()

			// TODO add after parsing as an independent query group, rather than as a
			// literal
			if defaultQueryGroup != "" {
				args = append(
					args,
					workspaceConfig.GetDefaultQueryGroup(),
				)
			}
		}
	}

	var latent errors.Multi

	if err, latent = state.build(args...); err != nil {
		if !errors.IsBadRequest(err) {
			if latent == nil {
				latent = errors.MakeMulti()
			}

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

	ui.Log().Print(state.qg.StringDebug())

	return
}
