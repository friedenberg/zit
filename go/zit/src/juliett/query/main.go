package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Query struct {
	kennung.Sigil
	kennung.Genre
	Exp

	Kennung map[string]Kennung

	Hidden sku.Query
}

func (a *Query) IsEmpty() bool {
	return a.Sigil == kennung.SigilUnknown &&
		a.Genre.IsEmpty() &&
		len(a.Children) == 0 &&
		len(a.Kennung) == 0
}

func (a *Query) GetSigil() kennung.Sigil {
	return a.Sigil
}

func (a *Query) ContainsKennung(k *kennung.Kennung2) bool {
	if !a.Genre.Contains(k.GetGenre()) {
		panic("should never check for wrong gattung")
	}

	if len(a.Kennung) == 0 {
		return false
	}

	_, ok := a.Kennung[k.String()]

	return ok
}

func (a *Query) Clone() (b *Query) {
	b = &Query{
		Sigil:   a.Sigil,
		Genre: a.Genre,
		Kennung: make(map[string]Kennung, len(a.Kennung)),
		Hidden:  a.Hidden,
	}

	bExp := a.Exp.Clone()
	b.Exp = *bExp

	for k, v := range a.Kennung {
		b.Kennung[k] = v
	}

	return b
}

func (q *Query) Add(m sku.Query) (err error) {
	q1, ok := m.(*Query)

	if !ok {
		return q.Exp.Add(m)
	}

	if q1.Genre != q.Genre {
		err = errors.Errorf(
			"expected %q but got %q",
			q.Genre,
			q1.Genre,
		)

		return
	}

	if err = q.Merge(q1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Query) Merge(b *Query) (err error) {
	a.Sigil.Add(b.Sigil)

	if a.Kennung == nil {
		a.Kennung = make(map[string]Kennung)
	}

	for _, k := range b.Kennung {
		a.Kennung[k.Kennung2.String()] = k
	}

	a.Children = append(a.Children, b.Children...)

	return
}

func (q *Query) Each(_ interfaces.FuncIter[sku.Query]) (err error) {
	return
}

func (q *Query) MatcherLen() int {
	return 0
}

func (q *Query) StringDebug() string {
	var sb strings.Builder

	if q.Kennung == nil || len(q.Kennung) == 0 {
		sb.WriteString(q.Exp.StringDebug())
	} else {
		sb.WriteString("[[")

		first := true

		for _, k := range q.Kennung {
			if !first {
				sb.WriteString(", ")
			}

			sb.WriteString(k.String())

			first = false
		}

		sb.WriteString(", ")
		sb.WriteString(q.Exp.StringDebug())
		sb.WriteString("]")
	}

	if q.IsEmpty() && !q.IsSchwanzenOrUnknown() {
		sb.WriteString(q.Sigil.String())
	} else if !q.IsEmpty() {
		sb.WriteString(q.Sigil.String())
		sb.WriteString(q.Genre.String())
	}

	return sb.String()
}

func (q *Query) SortedKennungen() []string {
	out := make([]string, 0, len(q.Kennung))

	for k := range q.Kennung {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

func (q *Query) String() string {
	var sb strings.Builder

	e := q.Exp.String()

	if q.Kennung == nil || len(q.Kennung) == 0 {
		sb.WriteString(e)
	} else if len(q.Kennung) == 1 && e == "" {
		for _, k := range q.Kennung {
			sb.WriteString(k.String())
		}
	} else {
		sb.WriteString("[")

		first := true

		for _, k := range q.SortedKennungen() {
			if !first {
				sb.WriteString(", ")
			}

			sb.WriteString(k)

			first = false
		}

		if e != "" {
			sb.WriteString(", ")
			sb.WriteString(q.Exp.String())
		}

		sb.WriteString("]")
	}

	if q.Genre.IsEmpty() && !q.IsSchwanzenOrUnknown() {
		sb.WriteString(q.Sigil.String())
	} else if !q.Genre.IsEmpty() {
		sb.WriteString(q.Sigil.String())
		sb.WriteString(q.Genre.String())
	}

	return sb.String()
}

func (q *Query) ShouldHide(sk *sku.Transacted, k string) bool {
	_, ok := q.Kennung[k]

	if q.IncludesHidden() || q.Hidden == nil || ok {
		return false
	}

	return q.Hidden.ContainsSku(sk)
}

func (q *Query) ContainsSku(sk *sku.Transacted) (ok bool) {
	defer sk.Metadatei.Verzeichnisse.QueryPath.PushOnOk(q, &ok)
	k := sk.Kennung.String()

	if q.ShouldHide(sk, k) {
		return
	}

	g := gattung.Must(sk)

	if !q.Genre.ContainsOneOf(g) {
		return
	}

	if _, ok = q.Kennung[k]; ok {
		return
	}

	if len(q.Children) == 0 {
		ok = len(q.Kennung) == 0 && q.MatchOnEmpty
		return
	} else if !q.Exp.ContainsSku(sk) {
		return
	}

	ok = true

	return
}
