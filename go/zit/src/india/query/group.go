package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func MakeQueryGroup(
	b *Builder,
) *QueryGroup {
	return &QueryGroup{
		OptimizedQueries: make(map[gattung.Gattung]*QueryWithHidden),
		UserQueries:      make(map[kennung.Gattung]*QueryWithHidden),
		Hidden:           b.hidden,
		FDs:              fd.MakeMutableSet(),
		Zettelen:         kennung.MakeHinweisMutableSet(),
		Etiketten:        kennung.MakeMutableEtikettSet(),
		Typen:            kennung.MakeMutableTypSet(),
	}
}

type QueryGroup struct {
	Hidden           matcher_proto.Matcher
	OptimizedQueries map[gattung.Gattung]*QueryWithHidden
	UserQueries      map[kennung.Gattung]*QueryWithHidden
	FDs              fd.MutableSet
	Zettelen         kennung.HinweisMutableSet
	Etiketten        kennung.EtikettMutableSet
	Typen            kennung.TypMutableSet
}

func (qg *QueryGroup) GetQueryGroup() matcher_proto.QueryGroup {
	return qg
}

func (qg *QueryGroup) Get(g gattung.Gattung) (matcher_proto.MatcherSigil, bool) {
	q, ok := qg.OptimizedQueries[g]
	return q, ok
}

func (qg *QueryGroup) GetCwdFDs() fd.Set {
	// TODO support dot operator
	// if ms.dotOperatorActive {
	// 	return ms.cwd.GetCwdFDs()
	// } else {
	// 	return ms.FDs
	// }
	return qg.FDs
}

func (qg *QueryGroup) GetExplicitCwdFDs() fd.Set {
	return qg.FDs
}

func (qg *QueryGroup) GetEtiketten() kennung.EtikettSet {
	return qg.Etiketten
}

func (qg *QueryGroup) GetTypen() kennung.TypSet {
	return qg.Typen
}

func (qg *QueryGroup) GetGattungen() gattungen.Set {
	gs := gattungen.MakeMutableSet()

	for g := range qg.OptimizedQueries {
		gs.Add(g)
	}

	return gs
}

func (qg *QueryGroup) AddExactKennung(
	b *Builder,
	k *kennung.Kennung2,
) (err error) {
	q := b.makeQuery()
	q.Kennung[k.String()] = k
	q.Gattung.Add(gattung.Must(k))

	if err = qg.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (qg *QueryGroup) Add(q *Query) (err error) {
	existing, ok := qg.UserQueries[q.Gattung]

	if !ok {
		existing = &QueryWithHidden{
			Hidden: qg.Hidden,
			Query: Query{
				Gattung: q.Gattung,
				Kennung: make(map[string]*kennung.Kennung2),
			},
		}
	}

	if err = existing.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.UserQueries[q.Gattung] = existing

	if err = qg.addOptimized(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (qg *QueryGroup) addOptimized(q *Query) (err error) {
	q = q.Clone()
	gs := q.Slice()

	for _, g := range gs {
		existing, ok := qg.OptimizedQueries[g]

		if !ok {
			existing = &QueryWithHidden{
				Hidden: qg.Hidden,
				Query: Query{
					Gattung: kennung.MakeGattung(g),
					Kennung: make(map[string]*kennung.Kennung2),
				},
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

func (q *QueryGroup) Each(_ schnittstellen.FuncIter[matcher_proto.Matcher]) (err error) {
	return
}

func (q *QueryGroup) MatcherLen() int {
	return 0
}

func (qg *QueryGroup) SortedUserQueries() []*QueryWithHidden {
	out := make([]*QueryWithHidden, 0, len(qg.UserQueries))

	for _, g := range qg.UserQueries {
		out = append(out, g)
	}

	sort.Slice(out, func(i, j int) bool {
		l, r := out[i].Gattung, out[j].Gattung

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

func (qg *QueryGroup) StringDebug() string {
	var sb strings.Builder

	first := true

	for _, g := range qg.SortedUserQueries() {
		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(g.StringDebug())

		first = false
	}

	return sb.String()
}

func (qg *QueryGroup) StringOptimized() string {
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

	for _, g := range qg.OptimizedQueries {
		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(g.String())

		first = false
	}

	return sb.String()
}

func (qg *QueryGroup) String() string {
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
		if !first {
			sb.WriteRune(' ')
		}

		sb.WriteString(g.String())

		first = false
	}

	return sb.String()
}

func (qg *QueryGroup) ContainsMatchable(sk *sku.Transacted) bool {
	log.Log().Print(qg, sk)
	g := sk.GetGattung()

	switch g {
	case gattung.Zettel:
		if qg.Zettelen.ContainsKey(sk.Kennung.String()) {
			return true
		}

	case gattung.Etikett:
		if qg.Etiketten.ContainsKey(sk.Kennung.String()) {
			return true
		}

	case gattung.Typ:
		if qg.Typen.ContainsKey(sk.Kennung.String()) {
			return true
		}
		// TODO other gattung
	}

	q, ok := qg.OptimizedQueries[gattung.Must(g)]

	if !ok {
		return false
	}

	if !q.ContainsMatchable(sk) {
		return false
	}

	return true
}
