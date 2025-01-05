package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/box"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/tag_blobs"
)

type stackEl interface {
	sku.Query
	Add(sku.Query) error
}

type buildState struct {
	builder      *Builder
	qg           *Group
	latentErrors errors.Multi
	missingBlobs []ErrBlobMissing

	luaVMPoolBuilder        *lua.VMPoolBuilder
	pinnedObjectIds         []pinnedObjectId
	pinnedExternalObjectIds []sku.ExternalObjectId
	repo                    sku.ExternalStoreForQuery
	eqo                     sku.ExternalQueryOptions

	externalStoreAcceptedQueryComponent bool

	scanner box.Scanner
}

func (b *buildState) makeGroup() *Group {
	return &Group{
		Hidden:           b.builder.hidden,
		OptimizedQueries: make(map[genres.Genre]*Query),
		UserQueries:      make(map[ids.Genre]*Query),
		Types:            ids.MakeMutableTypeSet(),
	}
}

func (b *buildState) build(
	vs ...string,
) (err error, latent errors.Multi) {
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
				if v != "." {
					remaining = append(remaining, v)
				}

				em.Add(err)
				err = nil

				continue
			}

			b.externalStoreAcceptedQueryComponent = true

			for _, k := range k {
				if k.GetGenre() == genres.None {
					err = errors.Errorf("id with empty genre: %q", k)
					return
				}

				b.pinnedExternalObjectIds = append(
					b.pinnedExternalObjectIds,
					k,
				)
			}
		}
	}

	remainingWithSpaces := make([]string, 0, len(remaining)*2)

	for i, s := range remaining {
		if i > 0 {
			remainingWithSpaces = append(remainingWithSpaces, " ")
		}

		remainingWithSpaces = append(remainingWithSpaces, s)
	}

	reader := catgut.MakeMultiRuneReader(remainingWithSpaces...)
	b.scanner.Reset(reader)

	for b.scanner.CanScan() {
		if err = b.parseTokens(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, k := range b.pinnedExternalObjectIds {
		if k.GetGenre() == genres.None {
			err = errors.Errorf("id with empty genre: %q", k)
			return
		}

		if err = b.qg.addExactExternalObjectId(b, k); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, k := range b.pinnedObjectIds {
		q := b.makeQuery()

		if err = q.addPinnedObjectId(b, k); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = b.qg.Add(q); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	b.addDefaultsIfNecessary()

	if err = b.qg.reduce(b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *buildState) addDefaultsIfNecessary() {
	if b.builder.defaultGenres.IsEmpty() || !b.qg.IsEmpty() {
		return
	}

	if b.builder.requireNonEmptyQuery && b.qg.IsEmpty() {
		return
	}

	if b.externalStoreAcceptedQueryComponent {
		return
	}

	b.qg.matchOnEmpty = true

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

func (b *buildState) parseTokens() (err error) {
	q := b.makeQuery()
	stack := []stackEl{q}

	isNegated := false
	isExact := false

LOOP:
	for b.scanner.Scan() {
		seq := b.scanner.GetSeq()

		if seq.MatchAll(box.TokenTypeOperator) {
			op := seq.At(0).Contents[0]

			switch op {
			case '=':
				isExact = true

			case '^':
				isNegated = true

			case ' ':
				if len(stack) == 1 {
					break LOOP
				}

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

				b.scanner.Unscan()

				if err = b.parseSigilsAndGenres(q); err != nil {
					err = errors.Wrapf(err, "Seq: %q", seq)
					return
				}

				continue LOOP
			}
		} else {
			if ok, left, right, partition := seq.PartitionFavoringRight(
				box.TokenMatcherOp(box.OpSigilExternal),
			); ok {
				switch {
				case right.MatchAll(box.TokenTypeIdentifier):
					if err = q.AddString(string(right.At(0).Contents)); err != nil {
						err = nil
					} else {
						if err = b.addSigilFromOp(q, partition.Contents[0]); err != nil {
							err = errors.Wrap(err)
							return
						}

						seq = left
					}

				case right.Len() == 0:
					if err = b.addSigilFromOp(q, partition.Contents[0]); err != nil {
						err = nil
					} else {
						seq = left
					}
				}
			}

			k := ObjectId{
				ObjectId: ids.GetObjectIdPool().Get(),
			}

			// TODO if this fails, permit an external store to try to read this as an
			// external object ID. And if that fails, try to remove the last two
			// elements as per the above and read that and force the genre and sigils
			if err = k.GetObjectId().ReadFromSeq(seq); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = k.reduce(b); err != nil {
				err = errors.Wrap(err)
				return
			}

			pid := pinnedObjectId{
				Sigil:    ids.SigilLatest,
				ObjectId: k,
			}

			switch k.GetGenre() {
			case genres.InventoryList, genres.Zettel:
				b.pinnedObjectIds = append(
					b.pinnedObjectIds,
					pid,
				)

				if err = q.addPinnedObjectId(b, pid); err != nil {
					err = errors.Wrap(err)
					return
				}

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

	if err = b.scanner.Error(); err != nil {
		err = errors.Wrap(err)
		return
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

func (b *buildState) addSigilFromOp(q *Query, op byte) (err error) {
	var s ids.Sigil

	if err = s.SetByte(op); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !b.builder.permittedSigil.IsEmpty() && !b.builder.permittedSigil.ContainsOneOf(s) {
		err = errors.BadRequestf("this query cannot contain the %q sigil", s)
		return
	}

	q.Sigil.Add(s)

	return
}

func (b *buildState) parseSigilsAndGenres(
	q *Query,
) (err error) {
	for b.scanner.Scan() {
		seq := b.scanner.GetSeq()

		if seq.MatchAll(box.TokenTypeOperator) {
			op := seq.At(0).Contents[0]

			switch op {
			default:
				err = errors.Errorf("unexpected operator %q", seq)
				return

			case ' ':
				return

			case '.':
				b.qg.dotOperatorActive = true
				fallthrough

			case ':', '+', '?':
				if err = b.addSigilFromOp(q, op); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		} else if seq.MatchAll(box.TokenTypeIdentifier) {
			b.scanner.Unscan()
			break
		} else {
			err = errors.Errorf("expected operator but got %q", seq)
			return
		}
	}

	if err = b.scanner.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = q.ReadFromBoxScanner(&b.scanner); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO use new generic and typed blobs
func (b *buildState) makeTagOrLuaTag(
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

	var twb sku.TransactedWithBlob[tag_blobs.Blob]

	if twb, _, err = b.builder.blob_store.GetTag().GetTransactedWithBlob(
		sk,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var matcherBlob sku.Queryable

	{
		var ok bool

		if matcherBlob, ok = twb.Blob.(sku.Queryable); !ok {
			return
		}
	}

	exp = &CompoundMatch{Queryable: matcherBlob, ObjectId: k}

	return
}

func (b *buildState) makeTagExp(k *ObjectId) (exp sku.Query, err error) {
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

func (b *buildState) makeExp(
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

func (b *buildState) makeQuery() *Query {
	return &Query{
		ObjectIds:         make(map[string]ObjectId),
		ExternalObjectIds: make(map[string]sku.ExternalObjectId),
	}
}
