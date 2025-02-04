package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type Group struct {
	Hidden           sku.Query
	OptimizedQueries map[genres.Genre]*Query
	UserQueries      map[ids.Genre]*Query
	Types            ids.TypeMutableSet

	dotOperatorActive bool
	matchOnEmpty      bool

	sku.ExternalQueryOptions
}

func (qg *Group) isDotOperatorActive() bool {
	if qg.dotOperatorActive {
		return true
	}

	for _, oq := range qg.OptimizedQueries {
		if oq.Sigil.ContainsOneOf(ids.SigilExternal) {
			return true
		}
	}

	return false
}

type reducer interface {
	reduce(*buildState) error
}

func (qg *Group) reduce(b *buildState) (err error) {
	for _, q := range qg.UserQueries {
		if err = q.reduce(b); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = qg.addOptimized(b, q); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, q := range qg.OptimizedQueries {
		if err = q.reduce(b); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (qg *Group) addExactExternalObjectId(
	b *buildState,
	k sku.ExternalObjectId,
) (err error) {
	if k == nil {
		err = errors.Errorf("nil object id")
		return
	}

	q := b.makeQuery()

	q.Sigil.Add(ids.SigilExternal)
	q.Sigil.Add(ids.SigilLatest)
	q.Genre.Add(genres.Must(k))
	q.ExternalObjectIds[k.String()] = k

	if err = qg.add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.dotOperatorActive = true

	return
}

func (qg *Group) add(q *Query) (err error) {
	existing, ok := qg.UserQueries[q.Genre]

	if !ok {
		existing = &Query{
			Hidden:    qg.Hidden,
			Genre:     q.Genre,
			ObjectIds: make(map[string]ObjectId),
		}
	}

	if err = existing.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.UserQueries[q.Genre] = existing

	return
}

func (qg *Group) addOptimized(b *buildState, q *Query) (err error) {
	q = q.Clone()
	gs := q.Slice()

	if len(gs) == 0 {
		gs = b.builder.defaultGenres.Slice()
	}

	for _, g := range gs {
		existing, ok := qg.OptimizedQueries[g]

		if !ok {
			existing = b.makeQuery()
			existing.Genre = ids.MakeGenre(g)
		}

		if err = existing.Merge(q); err != nil {
			err = errors.Wrap(err)
			return
		}

		qg.OptimizedQueries[g] = existing
	}

	return
}
