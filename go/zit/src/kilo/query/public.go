package query

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func IsExactlyOneObjectId(qg *Query) bool {
	if len(qg.optimizedQueries) != 1 {
		return false
	}

	var q *expSigilAndGenre

	for _, q1 := range qg.optimizedQueries {
		q = q1
	}

	kn := q.expObjectIds.internal
	lk := len(kn)

	if lk != 1 {
		return false
	}

	return true
}

func GetTags(query *Query) ids.TagSet {
	mes := ids.MakeMutableTagSet()

	for _, oq := range query.optimizedQueries {
		oq.expTagsOrTypes.CollectTags(mes)
	}

	return mes
}

func GetTypes(qg *Query) ids.TypeSet {
	return qg.types
}

func (qg *Query) String() string {
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

	for _, g := range qg.sortedUserQueries() {
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

func ContainsExternalSku(
	qg *Query,
	el sku.ExternalLike,
	state checked_out_state.State,
) (ok bool) {
	if qg.defaultQuery != nil &&
		!ContainsExternalSku(qg.defaultQuery, el, state) {
		return
	}

	sk := el.GetSku()

	defer sk.Metadata.Cache.QueryPath.PushOnReturn(qg, &ok)

	if !ContainsSkuCheckedOutState(qg, state) {
		return
	}

	if len(qg.optimizedQueries) == 0 && qg.matchOnEmpty {
		ok = true
		return
	}

	g := genres.Must(sk.GetGenre())

	q, ok := qg.optimizedQueries[g]

	if !ok || !q.ContainsExternalSku(el) {
		ok = false
		return
	}

	ok = true

	return
}

func ContainsSkuCheckedOutState(
	qg *Query,
	state checked_out_state.State,
) (ok bool) {
	if qg.defaultQuery != nil &&
		!ContainsSkuCheckedOutState(qg.defaultQuery, state) {
		return
	}

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
