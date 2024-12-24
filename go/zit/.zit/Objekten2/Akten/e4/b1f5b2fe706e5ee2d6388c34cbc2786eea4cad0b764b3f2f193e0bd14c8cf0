package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
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

func (qg *Group) DotOperatorActive() bool {
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

func (qg *Group) SetIncludeHistory() {
	for _, q := range qg.UserQueries {
		q.Sigil.Add(ids.SigilHistory)
	}

	for _, q := range qg.OptimizedQueries {
		q.Sigil.Add(ids.SigilHistory)
	}
}

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

	kn := q.ObjectIds
	lk := len(kn)

	if lk != 1 {
		return false
	}

	return true
}

func (qg *Group) GetExactlyOneObjectId(
	g genres.Genre,
) (k *ids.ObjectId, s ids.Sigil, err error) {
	if len(qg.OptimizedQueries) != 1 {
		err = errors.Errorf(
			"expected exactly 1 gattung query but got %d",
			len(qg.OptimizedQueries),
		)

		return
	}

	q, ok := qg.OptimizedQueries[g]

	if !ok {
		err = errors.Errorf("expected to have gattung %q", g)
		return
	}

	kn := q.ObjectIds
	lk := len(kn)

	if lk != 1 {
		err = errors.Errorf("expected to exactly 1 object id but got %d", lk)
		return
	}

	s = q.GetSigil()

	for _, k1 := range kn {
		k = k1.GetObjectId()

		// TODO
		// if k1.External {
		// 	s.Add(ids.SigilExternal)
		// }

		break
	}

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

func (qg *Group) GetGenres() (g ids.Genre) {
	for g1 := range qg.OptimizedQueries {
		g.Add(g1)
	}

	return
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

	if err = qg.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.dotOperatorActive = true

	return
}

func (qg *Group) Add(q *Query) (err error) {
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

func (q *Group) Each(_ interfaces.FuncIter[sku.Query]) (err error) {
	return
}

func (q *Group) MatcherLen() int {
	return 0
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

func (qg *Group) StringDebug() string {
	var sb strings.Builder

	first := true

	for _, g := range qg.SortedUserQueries() {
		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(g.StringDebug())

		first = false
	}

	sb.WriteString(" | ")
	first = true

	for _, g := range genres.TrueGenre() {
		q, ok := qg.OptimizedQueries[g]

		if !ok {
			continue
		}

		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(q.String())

		first = false
	}

	return sb.String()
}

func (qg *Group) StringOptimized() string {
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

	for _, g := range genres.TrueGenre() {
		q, ok := qg.OptimizedQueries[g]

		if !ok {
			continue
		}

		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(q.String())

		first = false
	}

	return sb.String()
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
		ui.Debug().Print(state)
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
