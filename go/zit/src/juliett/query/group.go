package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeGroup(
	b *Builder,
) *Group {
	return &Group{
		OptimizedQueries: make(map[genres.Genre]*Query),
		UserQueries:      make(map[ids.Genre]*Query),
		Hidden:           b.hidden,
		Zettels:          ids.MakeZettelIdMutableSet(),
		Types:            ids.MakeMutableTypeSet(),
	}
}

type Group struct {
	Hidden           sku.Query
	OptimizedQueries map[genres.Genre]*Query
	UserQueries      map[ids.Genre]*Query
	ObjectIds        []*ids.ObjectId
	Zettels          ids.ZettelIdMutableSet
	Types            ids.TypeMutableSet

	sku.ExternalQueryOptions
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
		err = errors.Errorf("expected to exactly 1 kennung but got %d", lk)
		return
	}

	s = q.GetSigil()

	for _, k1 := range kn {
		k = k1.ObjectId

		if k1.External {
			s.Add(ids.SigilExternal)
		}

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

type Reducer interface {
	Reduce(*Builder) error
}

func (qg *Group) Reduce(b *Builder) (err error) {
	for _, q := range qg.UserQueries {
		if err = q.Reduce(b); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = qg.addOptimized(b, q); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, q := range qg.OptimizedQueries {
		if err = q.Reduce(b); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (qg *Group) AddExactObjectId(
	b *Builder,
	k ObjectId,
) (err error) {
	if k.ObjectId == nil {
		err = errors.Errorf("nil kennung")
		return
	}

	qg.ObjectIds = append(qg.ObjectIds, k.ObjectId)

	q := b.makeQuery()
	q.Sigil.Add(ids.SigilLatest)
	q.ObjectIds[k.ObjectId.String()] = k
	q.Genre.Add(genres.Must(k))

	if err = qg.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (qg *Group) addOptimized(b *Builder, q *Query) (err error) {
	q = q.Clone()
	gs := q.Slice()

	if len(gs) == 0 {
		gs = b.defaultGenres.Slice()
	}

	for _, g := range gs {
		existing, ok := qg.OptimizedQueries[g]

		if !ok {
			existing = &Query{
				Hidden:    qg.Hidden,
				Genre:     ids.MakeGenre(g),
				ObjectIds: make(map[string]ObjectId),
			}
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

func (qg *Group) ContainsSku(sk *sku.Transacted) (ok bool) {
	defer sk.Metadata.Cache.QueryPath.PushOnOk(qg, &ok)
	g := sk.GetGenre()

	q, ok := qg.OptimizedQueries[genres.Must(g)]

	if !ok || !q.ContainsSku(sk) {
		ok = false
		return
	}

	ok = true

	return
}

func (qg *Group) MakeEmitSku(
	f interfaces.FuncIter[*sku.Transacted],
) interfaces.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := genres.Must(z.GetGenre())
		m, ok := qg.Get(g)

		if !ok {
			return
		}

		if !m.ContainsSku(z) {
			return
		}

		if err = f(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

// TODO improve performance by only reading Cwd zettels rather than scanning
// everything
func (qg *Group) MakeEmitSkuMaybeExternal(
	f interfaces.FuncIter[*sku.Transacted],
	k ids.RepoId,
	updateTransacted func(
		kasten ids.RepoId,
		z *sku.Transacted,
	) (err error),
) interfaces.FuncIter[*sku.Transacted] {
	// TODO add untracked and recognized
	// if qg.IncludeRecognized {
	// }

	// if !qg.ExcludeUntracked {
	// }

	if qg.GetSigil() == ids.SigilExternal {
		return qg.MakeEmitSkuSigilExternal(f, k, updateTransacted)
	} else {
		return qg.MakeEmitSkuSigilLatest(f, k, updateTransacted)
	}
}

func (qg *Group) MakeEmitSkuSigilLatest(
	f interfaces.FuncIter[*sku.Transacted],
	k ids.RepoId,
	updateTransacted func(
		kasten ids.RepoId,
		z *sku.Transacted,
	) (err error),
) interfaces.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := genres.Must(z.GetGenre())
		m, ok := qg.Get(g)

		if !ok {
			return
		}

		if m.GetSigil().IncludesExternal() {
			if err = updateTransacted(k, z); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !m.ContainsSku(z) {
			return
		}

		if err = f(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (qg *Group) MakeEmitSkuSigilExternal(
	f interfaces.FuncIter[*sku.Transacted],
	k ids.RepoId,
	updateTransacted func(
		kasten ids.RepoId,
		z *sku.Transacted,
	) (err error),
) interfaces.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := genres.Must(z.GetGenre())
		m, ok := qg.Get(g)

		if !ok {
			return
		}

		if err = updateTransacted(k, z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !m.ContainsSku(z) {
			return
		}

		if err = f(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
