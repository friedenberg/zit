package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Query struct {
	ids.Sigil
	ids.Genre
	Exp

	ObjectIds         map[string]ObjectId
	ExternalObjectIds map[string]sku.ExternalObjectId

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

func (q *Query) addPinnedObjectId(
	b *buildState,
	k pinnedObjectId,
) (err error) {
	if err = q.addExactObjectId(b, k.ObjectId, k.Sigil); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (q *Query) addExactObjectId(
	b *buildState,
	k ObjectId,
	sigil ids.Sigil,
) (err error) {
	if k.ObjectId == nil {
		err = errors.Errorf("nil object id")
		return
	}

	q.Sigil.Add(sigil)
	q.ObjectIds[k.GetObjectId().String()] = k
	q.Genre.Add(genres.Must(k))

	return
}

func (a *Query) ContainsObjectId(k *ids.ObjectId) bool {
	if !a.Genre.Contains(k.GetGenre()) {
		err := errors.Errorf("checking query %#v for object id %#v, %q, %q", a, k, a, k)
		panic(err)
	}

	if len(a.ObjectIds) == 0 {
		return false
	}

	_, ok := a.ObjectIds[k.String()]

	return ok
}

func (a *Query) Clone() (b *Query) {
	b = &Query{
		Sigil:             a.Sigil,
		Genre:             a.Genre,
		ObjectIds:         make(map[string]ObjectId, len(a.ObjectIds)),
		ExternalObjectIds: make(map[string]sku.ExternalObjectId, len(a.ExternalObjectIds)),
		Hidden:            a.Hidden,
	}

	bExp := a.Exp.Clone()
	b.Exp = *bExp

	for k, v := range a.ObjectIds {
		b.ObjectIds[k] = v
	}

	for k, v := range a.ExternalObjectIds {
		b.ExternalObjectIds[k] = v
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
		a.ObjectIds = make(map[string]ObjectId, len(b.ObjectIds))
	}

	for _, k := range b.ObjectIds {
		a.ObjectIds[k.GetObjectId().String()] = k
	}

	if a.ExternalObjectIds == nil {
		a.ExternalObjectIds = make(map[string]sku.ExternalObjectId, len(b.ExternalObjectIds))
	}

	for _, k := range b.ExternalObjectIds {
		a.ExternalObjectIds[k.GetExternalObjectId().String()] = k
	}

	a.Children = append(a.Children, b.Children...)

	return
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
	out := make([]string, 0, len(q.ObjectIds)+len(q.ExternalObjectIds))

	for k := range q.ObjectIds {
		out = append(out, k)
	}

	for k := range q.ExternalObjectIds {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

func (q *Query) String() string {
	var sb strings.Builder

	e := q.Exp.String()

	oids := q.SortedObjectIds()

	if len(oids) == 0 {
		sb.WriteString(e)
	} else if len(oids) == 1 && e == "" {
		for _, k := range oids {
			sb.WriteString(k)
		}
	} else {
		sb.WriteString("[")

		first := true

		for _, k := range oids {
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

func (q *Query) ShouldHide(tg sku.TransactedGetter, k string) bool {
	_, ok := q.ObjectIds[k]

	if q.IncludesHidden() || q.Hidden == nil || ok {
		return false
	}

	return q.Hidden.ContainsSku(tg)
}

func (q *Query) ContainsSku(tg sku.TransactedGetter) (ok bool) {
	sk := tg.GetSku()

	defer sk.Metadata.Cache.QueryPath.PushOnReturn(q, &ok)
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
		ok = len(q.ObjectIds) == 0
		return
	} else if !q.Exp.ContainsSku(tg) {
		return
	}

	ok = true

	return
}

func (q *Query) ContainsExternalSku(el sku.ExternalLike) (ok bool) {
	sk := el.GetSku()

	defer sk.Metadata.Cache.QueryPath.PushOnReturn(q, &ok)

	g := genres.Must(sk)

	if !q.Genre.ContainsOneOf(g) {
		return
	}

	k := sk.ObjectId.String()

	if q.ShouldHide(el, k) {
		return
	}

	eoid := el.GetExternalObjectId().String()
	if _, ok = q.ExternalObjectIds[eoid]; ok {
		return
	}

	if _, ok = q.ObjectIds[k]; ok {
		return
	}

	if len(q.Children) == 0 {
		ok = len(q.ObjectIds) == 0 && len(q.ExternalObjectIds) == 0
		return
	} else if !q.Exp.ContainsSku(el) {
		return
	}

	ok = true

	return
}
