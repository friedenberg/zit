package query

import (
	"sort"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

type Query struct {
	ids.Sigil
	ids.Genre

	exp

	Hidden sku.Query
}

func (a *Query) IsEmpty() bool {
	return a.Sigil == ids.SigilUnknown &&
		a.Genre.IsEmpty() &&
		a.exp.IsEmpty()
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
	q.expObjectIds.internal[k.GetObjectId().String()] = k
	q.Genre.Add(genres.Must(k))

	return
}

func (a *Query) ContainsObjectId(k *ids.ObjectId) bool {
	if !a.Genre.Contains(k.GetGenre()) {
		err := errors.Errorf("checking query %#v for object id %#v, %q, %q", a, k, a, k)
		panic(err)
	}

	if len(a.expObjectIds.internal) == 0 {
		return false
	}

	_, ok := a.expObjectIds.internal[k.String()]

	return ok
}

func (a *Query) Clone() (b *Query) {
	b = &Query{
		Sigil: a.Sigil,
		Genre: a.Genre,
		exp: exp{
			expObjectIds: expObjectIds{
				internal: make(map[string]ObjectId, len(a.expObjectIds.internal)),
				external: make(map[string]sku.ExternalObjectId, len(a.expObjectIds.external)),
			},
		},
		Hidden: a.Hidden,
	}

	bExp := a.expTagsOrTypes.Clone()
	b.expTagsOrTypes = *bExp

	for k, v := range a.expObjectIds.internal {
		b.expObjectIds.internal[k] = v
	}

	for k, v := range a.expObjectIds.external {
		b.expObjectIds.external[k] = v
	}

	return b
}

func (q *Query) Add(m sku.Query) (err error) {
	q1, ok := m.(*Query)

	if !ok {
		return q.expTagsOrTypes.Add(m)
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

	if a.expObjectIds.internal == nil {
		a.expObjectIds.internal = make(map[string]ObjectId, len(b.expObjectIds.internal))
	}

	for _, k := range b.expObjectIds.internal {
		a.expObjectIds.internal[k.GetObjectId().String()] = k
	}

	if a.expObjectIds.external == nil {
		a.expObjectIds.external = make(map[string]sku.ExternalObjectId, len(b.expObjectIds.external))
	}

	for _, k := range b.expObjectIds.external {
		a.expObjectIds.external[k.GetExternalObjectId().String()] = k
	}

	a.expTagsOrTypes.Children = append(a.expTagsOrTypes.Children, b.expTagsOrTypes.Children...)

	return
}

func (q *Query) StringDebug() string {
	var sb strings.Builder

	if q.expObjectIds.internal == nil || len(q.expObjectIds.internal) == 0 {
		sb.WriteString(q.expTagsOrTypes.StringDebug())
	} else {
		sb.WriteString("[[")

		first := true

		for _, k := range q.expObjectIds.internal {
			if !first {
				sb.WriteString(", ")
			}

			sb.WriteString(k.String())

			first = false
		}

		sb.WriteString(", ")
		sb.WriteString(q.expTagsOrTypes.StringDebug())
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
	out := make([]string, 0, q.expObjectIds.Len())

	for k := range q.expObjectIds.internal {
		out = append(out, k)
	}

	for k := range q.expObjectIds.external {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

func (q *Query) String() string {
	var sb strings.Builder

	e := q.expTagsOrTypes.String()

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
			sb.WriteString(q.expTagsOrTypes.String())
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
	_, ok := q.expObjectIds.internal[k]

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

	if _, ok = q.expObjectIds.internal[k]; ok {
		return
	}

	if len(q.expTagsOrTypes.Children) == 0 {
		ok = len(q.expObjectIds.internal) == 0
		return
	} else if !q.expTagsOrTypes.ContainsSku(tg) {
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
	ui.Log().Print(eoid, q.expObjectIds.external, q.expObjectIds.internal)

	if _, ok = q.expObjectIds.external[eoid]; ok {
		return
	}

	if _, ok = q.expObjectIds.external[k]; ok {
		return
	}

	if _, ok = q.expObjectIds.internal[k]; ok {
		return
	}

	if len(q.expTagsOrTypes.Children) == 0 {
		ok = q.expObjectIds.IsEmpty()
		return
	} else if !q.expTagsOrTypes.ContainsSku(el) {
		return
	}

	ok = true

	return
}
