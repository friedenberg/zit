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

type Query struct {
	ids.Sigil
	ids.Genre
	Exp

	ObjectIds map[string]ObjectId

	Hidden sku.Query
}

func (a *Query) IsEmpty() bool {
	return a.Sigil == ids.SigilUnknown &&
		a.Genre.IsEmpty() &&
		len(a.Children) == 0 &&
		len(a.ObjectIds) == 0
}

func (a *Query) GetSigil() ids.Sigil {
	return a.Sigil
}

func (a *Query) ContainsObjectId(k *ids.ObjectId) bool {
	if !a.Genre.Contains(k.GetGenre()) {
		panic("should never check for wrong gattung")
	}

	if len(a.ObjectIds) == 0 {
		return false
	}

	_, ok := a.ObjectIds[k.String()]

	return ok
}

func (a *Query) Clone() (b *Query) {
	b = &Query{
		Sigil:     a.Sigil,
		Genre:     a.Genre,
		ObjectIds: make(map[string]ObjectId, len(a.ObjectIds)),
		Hidden:    a.Hidden,
	}

	bExp := a.Exp.Clone()
	b.Exp = *bExp

	for k, v := range a.ObjectIds {
		b.ObjectIds[k] = v
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

	if a.ObjectIds == nil {
		a.ObjectIds = make(map[string]ObjectId)
	}

	for _, k := range b.ObjectIds {
		a.ObjectIds[k.ObjectId.String()] = k
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

	if q.ObjectIds == nil || len(q.ObjectIds) == 0 {
		sb.WriteString(q.Exp.StringDebug())
	} else {
		sb.WriteString("[[")

		first := true

		for _, k := range q.ObjectIds {
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

	if q.IsEmpty() && !q.IsLatestOrUnknown() {
		sb.WriteString(q.Sigil.String())
	} else if !q.IsEmpty() {
		sb.WriteString(q.Sigil.String())
		sb.WriteString(q.Genre.String())
	}

	return sb.String()
}

func (q *Query) SortedObjectIds() []string {
	out := make([]string, 0, len(q.ObjectIds))

	for k := range q.ObjectIds {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

func (q *Query) String() string {
	var sb strings.Builder

	e := q.Exp.String()

	if q.ObjectIds == nil || len(q.ObjectIds) == 0 {
		sb.WriteString(e)
	} else if len(q.ObjectIds) == 1 && e == "" {
		for _, k := range q.ObjectIds {
			sb.WriteString(k.String())
		}
	} else {
		sb.WriteString("[")

		first := true

		for _, k := range q.SortedObjectIds() {
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

	if q.Genre.IsEmpty() && !q.IsLatestOrUnknown() {
		sb.WriteString(q.Sigil.String())
	} else if !q.Genre.IsEmpty() {
		sb.WriteString(q.Sigil.String())
		sb.WriteString(q.Genre.String())
	}

	return sb.String()
}

func (q *Query) ShouldHide(sk *sku.Transacted, k string) bool {
	_, ok := q.ObjectIds[k]

	if q.IncludesHidden() || q.Hidden == nil || ok {
		return false
	}

	return q.Hidden.ContainsSku(sk)
}

func (q *Query) ContainsSku(sk *sku.Transacted) (ok bool) {
	defer sk.Metadata.Cache.QueryPath.PushOnOk(q, &ok)
	k := sk.ObjectId.String()

	if q.ShouldHide(sk, k) {
		return
	}

	g := genres.Must(sk)

	if !q.Genre.ContainsOneOf(g) {
		return
	}

	if _, ok = q.ObjectIds[k]; ok {
		return
	}

	if len(q.Children) == 0 {
		ok = len(q.ObjectIds) == 0 && q.MatchOnEmpty
		return
	} else if !q.Exp.ContainsSku(sk) {
		return
	}

	ok = true

	return
}
