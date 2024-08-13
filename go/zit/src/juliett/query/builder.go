package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/tag_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
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
	repoGetter                 sku.ExternalStoreForQueryGetter
	repo                       sku.ExternalStoreForQuery
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

type builderState struct {
	builder      *Builder
	qg           *Group
	latentErrors errors.Multi

	luaVMPoolBuilder *lua.VMPoolBuilder
	pinnedObjectIds  []ObjectId
	repo             sku.ExternalStoreForQuery
	virtualEtiketten map[string]Lua
	eqo              sku.ExternalQueryOptions
}

func (b *Builder) makeState() *builderState {
	state := &builderState{
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

func (mb *Builder) WithRepo(
	repo sku.ExternalStoreForQuery,
) *Builder {
	mb.repo = repo
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
			b.pinnedObjectIds = append(
				b.pinnedObjectIds,
				ObjectId{
					ObjectIdLike: t.GetSku().ObjectId.Clone(),
				},
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

func (b *Builder) build(state *builderState, vs ...string) (err error) {
	var latent errors.Multi

	if err, latent = state.build(vs...); err != nil {
		latent.Add(errors.Wrapf(err, "Query: %q", vs))
		err = latent
		return
	}

	ui.Log().Print(state.qg.StringDebug())

	return
}

//   ____        _ _     _ _
//  | __ ) _   _(_) | __| (_)_ __   __ _
//  |  _ \| | | | | |/ _` | | '_ \ / _` |
//  | |_) | |_| | | | (_| | | | | | (_| |
//  |____/ \__,_|_|_|\__,_|_|_| |_|\__, |
//                                 |___/

func (b *builderState) makeGroup() *Group {
	return &Group{
		Hidden:            b.builder.hidden,
		OptimizedQueries:  make(map[genres.Genre]*Query),
		UserQueries:       make(map[ids.Genre]*Query),
		ExternalObjectIds: make(map[string]ObjectId),
		Types:             ids.MakeMutableTypeSet(),
	}
}

func (b *builderState) build(
	vs ...string,
) (err error, latent errors.Multi) {
	if err = b.realizeVirtualTags(); err != nil {
		err = errors.Wrap(err)
		return
	}

	em := errors.MakeMulti()
	latent = em

	var remaining []string

	if b.repo == nil {
		remaining = vs
	} else {
		for _, v := range vs {
			if v == "." {
				b.qg.dotOperatorActive = true
				remaining = append(remaining, v)
			}

			var k []sku.ExternalObjectId

			if k, err = b.repo.GetObjectIdsForString(v); err != nil {
				em.Add(err)
				err = nil
				remaining = append(remaining, v)
				continue
			}

			for _, k := range k {
				// b.defaultGenres.Add(genres.Must(k.GetObjectId().GetGenre()))

				b.pinnedObjectIds = append(
					b.pinnedObjectIds,
					ObjectId{
						ObjectIdLike: k,
						External:     true,
					},
				)
			}
		}
	}

	var tokens []string

	if tokens, err = query_spec.GetTokensFromStrings(remaining...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = b.buildManyFromTokens(tokens...); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, k := range b.pinnedObjectIds {
		if err = b.qg.addExactObjectId(b, k); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	b.addDefaultsIfNecessary()

	if err = b.qg.Reduce(b.builder); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *builderState) realizeVirtualTags() (err error) {
	for k, v := range b.builder.virtualEtikettenBeforeInit {
		var vmp *lua.VMPool

		lb := b.luaVMPoolBuilder.Clone().WithScript(v)

		if vmp, err = lb.Build(); err != nil {
			err = errors.Wrap(err)
			return
		}

		ml := Lua{
			LuaVMPool: sku_fmt.MakeLuaVMPool(vmp, nil),
		}

		b.virtualEtiketten[k] = ml
	}

	return
}

func (b *builderState) buildManyFromTokens(
	tokens ...string,
) (err error) {
	for len(tokens) > 0 {
		if tokens, err = b.parseOneFromTokens(tokens...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (b *builderState) addDefaultsIfNecessary() {
	// defer b.addDotOperatorIfNecessary()

	if b.builder.defaultGenres.IsEmpty() || !b.qg.IsEmpty() {
		return
	}

	if b.builder.requireNonEmptyQuery && b.qg.IsEmpty() {
		return
	}

	g := ids.MakeGenre()
	dq, ok := b.qg.UserQueries[g]

	if ok {
		delete(b.qg.UserQueries, g)
	} else {
		dq = b.makeQuery()
	}

	dq.Genre = b.builder.defaultGenres

	if b.builder.defaultSigil.IsEmpty() {
		dq.Sigil = ids.SigilLatest
	} else {
		dq.Sigil = b.builder.defaultSigil
	}

	b.qg.UserQueries[b.builder.defaultGenres] = dq
}

func (b *builderState) addDotOperatorIfNecessary() {
	if b.qg.dotOperatorActive {
		return
	}

	permitted := false

	for _, q := range b.qg.UserQueries {
		if q.Sigil.IncludesExternal() {
			permitted = true
			break
		}
	}

	if !permitted {
		return
	}

	b.qg.dotOperatorActive = true

	var k []sku.ExternalObjectId
	var err error

	if k, err = b.repo.GetObjectIdsForString("."); err != nil {
		b.latentErrors.Add(err)
		err = nil
	}

	for _, k := range k {
		// b.defaultGenres.Add(genres.Must(k.GetObjectId().GetGenre()))
		if err = b.qg.addExactObjectId(
			b,
			ObjectId{
				ObjectIdLike: k,
				External:     true,
			},
		); err != nil {
			b.latentErrors.Add(err)
			err = nil
		}
	}
}

func (b *builderState) makeQuery() *Query {
	return &Query{
		ObjectIds: make(map[string]ObjectId),
	}
}

func (b *builderState) makeExp(
	negated, exact bool,
	children ...sku.Query,
) *Exp {
	return &Exp{
		// MatchOnEmpty: !b.doNotMatchEmpty,
		Negated:  negated,
		Exact:    exact,
		Children: children,
	}
}

func (b *builderState) parseOneFromTokens(
	tokens ...string,
) (remainingTokens []string, err error) {
	type stackEl interface {
		sku.Query
		Add(sku.Query) error
	}

	q := b.makeQuery()
	stack := []stackEl{q}

	isNegated := false
	isExact := false

LOOP:
	for i, el := range tokens {
		// TODO refactor into separate functions
		if len(el) == 1 && query_spec.IsMatcherOperator([]rune(el)[0]) {
			op := el[0]
			switch op {
			case '=':
				isExact = true

			case '^':
				isNegated = true

			case ' ':

			case ',':
				last := stack[len(stack)-1].(*Exp)
				last.Or = true
				// TODO handle or when invalid

			case '[':
				exp := b.makeExp(isNegated, isExact)
				isExact = false
				isNegated = false
				stack[len(stack)-1].Add(exp)
				stack = append(stack, exp)

			case ']':
				stack = stack[:len(stack)-1]
				// TODO handle errors of unbalanced

			case '.':
				// TODO end sigil or embedded as part of name
				fallthrough

			case ':', '+', '?':
				if len(stack) > 1 {
					err = errors.Errorf("sigil before end")
					return
				}

				if remainingTokens, err = b.parseSigilsAndGenres(
					q,
					tokens[i:]...,
				); err != nil {
					err = errors.Wrapf(err, "%s", tokens[i:])
					return
				}

				break LOOP
			}
		} else {
			k := ObjectId{
				ObjectIdLike: ids.GetObjectIdPool().Get(),
			}

			if err = k.GetObjectId().Set(el); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = k.Reduce(b.builder); err != nil {
				err = errors.Wrap(err)
				return
			}

			switch k.GetGenre() {
			case genres.Zettel:
				b.pinnedObjectIds = append(
					b.pinnedObjectIds,
					k,
				)

				q.Genre.Add(genres.Zettel)
				q.ObjectIds[k.ObjectIdLike.String()] = k

			case genres.Tag:
				var et sku.Query

				if et, err = b.makeTagExp(&k); err != nil {
					err = errors.Wrap(err)
					return
				}

				exp := b.makeExp(isNegated, isExact, et)
				stack[len(stack)-1].Add(exp)

			case genres.Type:
				var t ids.Type

				if err = t.TodoSetFromObjectId(k.GetObjectId()); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !isNegated {
					if err = b.qg.Types.Add(t); err != nil {
						err = errors.Wrap(err)
						return
					}
				}

				exp := b.makeExp(isNegated, isExact, &k)
				stack[len(stack)-1].Add(exp)
			}

			isNegated = false
			isExact = false
		}
	}

	if q.IsEmpty() {
		return
	}

	if q.Genre.IsEmpty() && !b.builder.requireNonEmptyQuery {
		q.Genre = b.builder.defaultGenres
	}

	if q.Sigil.IsEmpty() {
		q.Sigil = b.builder.defaultSigil
	}

	if err = b.qg.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *builderState) makeTagOrLuaTag(
	k *ObjectId,
) (exp sku.Query, err error) {
	exp = k

	if b.builder.object_probe_index == nil || b.builder.blob_store == nil {
		return
	}

	sk := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(sk)

	if err = b.builder.object_probe_index.ReadOneObjectId(
		k.String(),
		sk,
	); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	lb := b.luaVMPoolBuilder.Clone().WithApply(MakeSelfApply(sk))

	// TODO use repo pattern
	if sk.GetType().String() == "lua" {
		var ar sha.ReadCloser

		if ar, err = b.builder.fs_home.BlobReader(sk.GetBlobSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, ar)

		lb.WithReader(ar)
	} else {
		var blob *tag_blobs.V1

		if blob, err = b.builder.blob_store.GetTagV1().GetBlob(
			sk.GetBlobSha(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if blob.Filter == "" {
			return
		}

		lb.WithScript(blob.Filter)
	}

	var vmp *lua.VMPool

	if vmp, err = lb.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ml := Lua{
		LuaVMPool: sku_fmt.MakeLuaVMPool(vmp, nil),
	}

	exp = &TagLua{Lua: &ml, ObjectId: k}

	return
}

func (b *builderState) makeTagExp(k *ObjectId) (exp sku.Query, err error) {
	// TODO use b.blobs to read tag blob and find filter if necessary
	var e ids.Tag

	if err = e.TodoSetFromObjectId(k.GetObjectId()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if exp, err = b.makeTagOrLuaTag(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *builderState) parseSigilsAndGenres(
	q *Query,
	tokens ...string,
) (remainingTokens []string, err error) {
LOOP:
	for i, el := range tokens {
		if len(el) != 1 {
			remainingTokens = tokens[i:]
			break
		}

		op := []rune(el)[0]

		switch op {
		default:
			remainingTokens = tokens[i:]
			break LOOP

		case ':', '+', '?', '.':
			var s ids.Sigil

			if err = s.Set(el); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !b.builder.permittedSigil.IsEmpty() && !b.builder.permittedSigil.ContainsOneOf(s) {
				err = errors.Errorf("cannot contain sigil %s", s)
				return
			}

			q.Sigil.Add(s)
		}
	}

	if remainingTokens, err = q.SetTokens(remainingTokens...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
