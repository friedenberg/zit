package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func MakeGroup(
	b *Builder,
) *Group {
	return &Group{
		OptimizedQueries: make(map[gattung.Gattung]*Query),
		UserQueries:      make(map[kennung.Gattung]*Query),
		Hidden:           b.hidden,
		Zettelen:         kennung.MakeHinweisMutableSet(),
		Typen:            kennung.MakeMutableTypSet(),
	}
}

type Group struct {
	Hidden           sku.Query
	OptimizedQueries map[gattung.Gattung]*Query
	UserQueries      map[kennung.Gattung]*Query
	Kennungen        []*kennung.Kennung2
	Zettelen         kennung.HinweisMutableSet
	Typen            kennung.TypMutableSet

	sku.ExternalQueryOptions
}

func (qg *Group) SetIncludeHistory() {
	for _, q := range qg.UserQueries {
		q.Sigil.Add(kennung.SigilHistory)
	}

	for _, q := range qg.OptimizedQueries {
		q.Sigil.Add(kennung.SigilHistory)
	}
}

func (qg *Group) HasHidden() bool {
	return qg.Hidden != nil
}

func (qg *Group) IsEmpty() bool {
	return len(qg.UserQueries) == 0
}

func (qg *Group) Get(g gattung.Gattung) (sku.QueryWithSigilAndKennung, bool) {
	q, ok := qg.OptimizedQueries[g]
	return q, ok
}

func (qg *Group) GetSigil() (s kennung.Sigil) {
	for _, q := range qg.OptimizedQueries {
		s.Add(q.Sigil)
	}

	return
}

func (qg *Group) GetExactlyOneKennung(
	g gattung.Gattung,
) (k *kennung.Kennung2, s kennung.Sigil, err error) {
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

	kn := q.Kennung
	lk := len(kn)

	if lk != 1 {
		err = errors.Errorf("expected to exactly 1 kennung but got %d", lk)
		return
	}

	s = q.GetSigil()

	for _, k1 := range kn {
		k = k1.Kennung2

		if k1.External {
			s.Add(kennung.SigilExternal)
		}

		break
	}

	return
}

func (qg *Group) GetEtiketten() kennung.EtikettSet {
	mes := kennung.MakeMutableEtikettSet()

	for _, oq := range qg.OptimizedQueries {
		oq.CollectEtiketten(mes)
	}

	return mes
}

func (qg *Group) GetTypen() kennung.TypSet {
	return qg.Typen
}

func (qg *Group) GetGattungen() (g kennung.Gattung) {
	for g1 := range qg.OptimizedQueries {
		g.Add(g1)
	}

	return
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
	k Kennung,
) (err error) {
	if k.Kennung2 == nil {
		err = errors.Errorf("nil kennung")
		return
	}

	qg.Kennungen = append(qg.Kennungen, k.Kennung2)

	q := b.makeQuery()
	q.Sigil.Add(kennung.SigilSchwanzen)
	q.Kennung[k.Kennung2.String()] = k
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
		existing = &Query{
			Hidden:  qg.Hidden,
			Gattung: q.Gattung,
			Kennung: make(map[string]Kennung),
		}
	}

	if err = existing.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.UserQueries[q.Gattung] = existing

	return
}

func (qg *Group) addOptimized(b *Builder, q *Query) (err error) {
	q = q.Clone()
	gs := q.Slice()

	if len(gs) == 0 {
		gs = b.defaultGattungen.Slice()
	}

	for _, g := range gs {
		existing, ok := qg.OptimizedQueries[g]

		if !ok {
			existing = &Query{
				Hidden:  qg.Hidden,
				Gattung: kennung.MakeGattung(g),
				Kennung: make(map[string]Kennung),
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

func (q *Group) Each(_ schnittstellen.FuncIter[sku.Query]) (err error) {
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

func (qg *Group) ContainsSku(sk *sku.Transacted) (ok bool) {
	defer sk.Metadatei.Verzeichnisse.QueryPath.PushOnOk(qg, &ok)
	g := sk.GetGattung()

	q, ok := qg.OptimizedQueries[gattung.Must(g)]

	if !ok || !q.ContainsSku(sk) {
		ok = false
		return
	}

	ok = true

	return
}

func (qg *Group) MakeEmitSku(
	f schnittstellen.FuncIter[*sku.Transacted],
) schnittstellen.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
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
	f schnittstellen.FuncIter[*sku.Transacted],
	k kennung.Kasten,
	updateTransacted func(
		kasten kennung.Kasten,
		z *sku.Transacted,
	) (err error),
) schnittstellen.FuncIter[*sku.Transacted] {
	// TODO add untracked and recognized
	// if qg.IncludeRecognized {
	// }

	// if !qg.ExcludeUntracked {
	// }

	if qg.GetSigil() == kennung.SigilExternal {
		return qg.MakeEmitSkuSigilExternal(f, k, updateTransacted)
	} else {
		return qg.MakeEmitSkuSigilSchwanzen(f, k, updateTransacted)
	}
}

func (qg *Group) MakeEmitSkuSigilSchwanzen(
	f schnittstellen.FuncIter[*sku.Transacted],
	k kennung.Kasten,
	updateTransacted func(
		kasten kennung.Kasten,
		z *sku.Transacted,
	) (err error),
) schnittstellen.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
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
	f schnittstellen.FuncIter[*sku.Transacted],
	k kennung.Kasten,
	updateTransacted func(
		kasten kennung.Kasten,
		z *sku.Transacted,
	) (err error),
) schnittstellen.FuncIter[*sku.Transacted] {
	return func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
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
