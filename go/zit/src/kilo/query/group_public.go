package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (qg *Group) HasHidden() bool {
	return qg.Hidden != nil
}

func (qg *Group) IsEmpty() bool {
	return len(qg.UserQueries) == 0
}

// TODO migrate this to the query executor
func (qg *Group) Get(g genres.Genre) (sku.QueryWithSigilAndObjectId, bool) {
	q, ok := qg.OptimizedQueries[g]
	return q, ok
}

func (qg *Group) GetSigil() (s ids.Sigil) {
	for _, q := range qg.OptimizedQueries {
		s.Add(q.Sigil)
	}

	return
}

func (qg *Group) IsExactlyOneObjectId() bool {
	if len(qg.OptimizedQueries) != 1 {
		return false
	}

	var q *Query

	for _, q1 := range qg.OptimizedQueries {
		q = q1
	}

	kn := q.internal
	lk := len(kn)

	if lk != 1 {
		return false
	}

	return true
}

func (queryGroup *Group) GetExactlyOneExternalObjectId(
	genre genres.Genre,
) (objectId ids.ExternalObjectIdLike, sigil ids.Sigil, err error) {
	if len(queryGroup.OptimizedQueries) != 1 {
		err = errors.Errorf(
			"expected exactly 1 genre query but got %d",
			len(queryGroup.OptimizedQueries),
		)

		return
	}

	query, ok := queryGroup.OptimizedQueries[genre]

	if !ok {
		err = errors.Errorf("expected to have genre %q", genre)
		return
	}

	if query.Sigil.ContainsOneOf(ids.SigilHistory) {
		err = errors.Errorf(
			"sigil (%s) includes history, which may return multiple objects",
			query.Sigil,
		)

		return
	}

	oids := query.internal
	oidsLen := len(oids)

	eoids := query.external
	eoidsLen := len(eoids)

	switch {
	case eoidsLen == 1 && oidsLen == 0:
		for _, k1 := range eoids {
			objectId = k1.GetExternalObjectId()
		}

		sigil.Add(ids.SigilExternal)

	default:
		err = errors.Errorf(
			"expected to exactly 1 object id or 1 external object id but got %d object ids and %d external object ids",
			oidsLen,
			eoidsLen,
		)

		return
	}

	sigil = query.GetSigil()

	return
}

func (queryGroup *Group) GetExactlyOneObjectId(
	genre genres.Genre,
) (objectId *ids.ObjectId, sigil ids.Sigil, err error) {
	if len(queryGroup.OptimizedQueries) != 1 {
		err = errors.Errorf(
			"expected exactly 1 genre query but got %d",
			len(queryGroup.OptimizedQueries),
		)

		return
	}

	query, ok := queryGroup.OptimizedQueries[genre]

	if !ok {
		err = errors.Errorf("expected to have genre %q", genre)
		return
	}

	if query.Sigil.ContainsOneOf(ids.SigilHistory) {
		err = errors.Errorf(
			"sigil (%s) includes history, which may return multiple objects",
			query.Sigil,
		)

		return
	}

	oids := query.internal
	oidsLen := len(oids)

	eoids := query.external
	eoidsLen := len(eoids)

	switch {
	case eoidsLen == 0 && oidsLen == 1:
		for _, k1 := range oids {
			objectId = k1.GetObjectId()
		}

	default:
		err = errors.Errorf(
			"expected to exactly 1 object id or 1 external object id but got %d object ids and %d external object ids",
			oidsLen,
			eoidsLen,
		)

		return
	}

	sigil = query.GetSigil()

	return
}

func (qg *Group) GetTags() ids.TagSet {
	mes := ids.MakeMutableTagSet()

	for _, oq := range qg.OptimizedQueries {
		oq.CollectTags(mes)
	}

	return mes
}

func (qg *Group) GetTypes() ids.TypeSet {
	return qg.Types
}

func (qg *Group) SortedUserQueries() []*Query {
	out := make([]*Query, 0, len(qg.UserQueries))

	for _, g := range qg.UserQueries {
		out = append(out, g)
	}

	sort.Slice(out, func(i, j int) bool {
		l, r := out[i].Genre, out[j].Genre

		if l.IsEmpty() {
			return false
		}

		if r.IsEmpty() {
			return true
		}

		return l < r
	})

	return out
}

func (qg *Group) String() string {
	var sb strings.Builder

	first := true

	// qg.FDs.Each(
	// 	func(f *fd.FD) error {
	// 		if !first {
	// 			sb.WriteRune(' ')
	// 		}

	// 		sb.WriteString(f.String())

	// 		first = false

	// 		return nil
	// 	},
	// )

	for _, g := range qg.SortedUserQueries() {
		// TODO determine why GS can be ""
		gs := g.String()

		if gs == "" {
			continue
		}

		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(gs)

		first = false
	}

	return sb.String()
}

func (qg *Group) ContainsSku(tg sku.TransactedGetter) (ok bool) {
	sk := tg.GetSku()
	defer sk.Metadata.Cache.QueryPath.PushOnReturn(qg, &ok)

	if len(qg.OptimizedQueries) == 0 && qg.matchOnEmpty {
		ok = true
		return
	}

	g := sk.GetGenre()

	q, ok := qg.OptimizedQueries[genres.Must(g)]

	if !ok || !q.ContainsSku(tg) {
		ok = false
		return
	}

	ok = true

	return
}

func (qg *Group) ContainsExternalSku(
	el sku.ExternalLike,
	state checked_out_state.State,
) (ok bool) {
	sk := el.GetSku()

	defer sk.Metadata.Cache.QueryPath.PushOnReturn(qg, &ok)

	if !qg.ContainsSkuCheckedOutState(state) {
		return
	}

	if len(qg.OptimizedQueries) == 0 && qg.matchOnEmpty {
		ok = true
		return
	}

	g := genres.Must(sk.GetGenre())

	q, ok := qg.OptimizedQueries[g]

	if !ok || !q.ContainsExternalSku(el) {
		ok = false
		return
	}

	ok = true

	return
}

func (qg *Group) ContainsSkuCheckedOutState(
	state checked_out_state.State,
) (ok bool) {
	switch state {
	case checked_out_state.Untracked:
		ok = !qg.ExcludeUntracked

	case checked_out_state.Recognized:
		ok = !qg.ExcludeRecognized

	default:
		ok = true
	}

	return
}
