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
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func MakeGroup(
	b *Builder,
) *Group {
	return &Group{
		OptimizedQueries: make(map[gattung.Gattung]*QueryWithHidden),
		UserQueries:      make(map[kennung.Gattung]*QueryWithHidden),
		Hidden:           b.hidden,
		FDs:              fd.MakeMutableSet(),
		Zettelen:         kennung.MakeHinweisMutableSet(),
		Etiketten:        kennung.MakeMutableEtikettSet(),
		Typen:            kennung.MakeMutableTypSet(),
	}
}

type Group struct {
	Hidden           Matcher
	OptimizedQueries map[gattung.Gattung]*QueryWithHidden
	UserQueries      map[kennung.Gattung]*QueryWithHidden
	FDs              fd.MutableSet
	Zettelen         kennung.HinweisMutableSet
	Etiketten        kennung.EtikettMutableSet
	Typen            kennung.TypMutableSet
}

func (qg *Group) IsEmpty() bool {
	return len(qg.UserQueries) == 0
}

func (qg *Group) Get(g gattung.Gattung) (MatcherSigil, bool) {
	q, ok := qg.OptimizedQueries[g]
	return q, ok
}

func (qg *Group) GetCwdFDs() fd.Set {
	// TODO support dot operator
	// if ms.dotOperatorActive {
	// 	return ms.cwd.GetCwdFDs()
	// } else {
	// 	return ms.FDs
	// }
	return qg.FDs
}

func (qg *Group) GetExplicitCwdFDs() fd.Set {
	return qg.FDs
}

func (qg *Group) GetEtiketten() kennung.EtikettSet {
	return qg.Etiketten
}

func (qg *Group) GetTypen() kennung.TypSet {
	return qg.Typen
}

func (qg *Group) GetGattungen() gattungen.Set {
	gs := gattungen.MakeMutableSet()

	for g := range qg.OptimizedQueries {
		gs.Add(g)
	}

	return gs
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

func (qg *Group) AddExactKennung(
	b *Builder,
	k *kennung.Kennung2,
) (err error) {
	q := b.makeQuery()
	q.Sigil.Add(kennung.SigilSchwanzen)
	q.Kennung[k.String()] = k
	q.Gattung.Add(gattung.Must(k))

	if err = qg.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (qg *Group) Add(q *Query) (err error) {
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

	return
}

func (qg *Group) addOptimized(b *Builder, q *QueryWithHidden) (err error) {
	q = q.Clone()
	gs := q.Slice()

	if len(gs) == 0 {
		gs = b.defaultGattungen.Slice()
	}

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

		if err = existing.Merge(&q.Query); err != nil {
			err = errors.Wrap(err)
			return
		}

		qg.OptimizedQueries[g] = existing
	}

	return
}

func (q *Group) Each(_ schnittstellen.FuncIter[Matcher]) (err error) {
	return
}

func (q *Group) MatcherLen() int {
	return 0
}

func (qg *Group) SortedUserQueries() []*QueryWithHidden {
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

	for _, g := range gattung.TrueGattung() {
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

	for _, g := range gattung.TrueGattung() {
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

func (qg *Group) ContainsMatchable(sk *sku.Transacted) bool {
	log.Log().Printf("%s in %s", sk, qg)
	g := sk.GetGattung()

	// switch g {
	// case gattung.Zettel:
	// 	if qg.Zettelen.ContainsKey(sk.Kennung.String()) {
	// 		return true
	// 	}

	// case gattung.Etikett:
	// 	if qg.Etiketten.ContainsKey(sk.Kennung.String()) {
	// 		return true
	// 	}

	// case gattung.Typ:
	// 	if qg.Typen.ContainsKey(sk.Kennung.String()) {
	// 		return true
	// 	}
	// 	// TODO other gattung
	// }

	q, ok := qg.OptimizedQueries[gattung.Must(g)]

	if !ok {
		return false
	}

	if !q.ContainsMatchable(sk) {
		return false
	}

	return true
}
