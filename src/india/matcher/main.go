package matcher

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func MakeMatcher(
	k kennung.KennungSansGattungPtr,
	v string,
	expander func(string) (string, error),
	ki kennung.Index,
	konfig schnittstellen.Konfig,
) (m Matcher, isNegated bool, isExact bool, err error) {
	v = strings.TrimSpace(v)
	didExpand := false

	if expander != nil {
		v1 := v

		if v1, err = expander(v); err != nil {
			err = nil
			v1 = v
		} else {
			didExpand = true
		}

		v = v1
	}

	if isNegated, isExact, err = SetQueryKennung(k, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !didExpand && expander != nil {
		v1 := k.String()

		if v1, err = expander(v1); err != nil {
			err = nil
			v1 = v
		}

		if err = k.Set(v1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if isExact {
		m = MakeMatcherContainsExactly(k)
	} else {
		m = MakeMatcherContains(k, ki)
	}

	if isNegated {
		m = MakeMatcherNegate(m)
	}

	// m = MakeMatcherAnd(
	// 	MakeMatcherLua(
	// 		`function contains_matchable(sku) return true end`,
	// 	),
	// 	m,
	// )

	return
}

func SetQueryKennung(
	k schnittstellen.Setter,
	v string,
) (isNegated bool, isExact bool, err error) {
	v = strings.TrimSpace(v)

	if len(v) > 0 && []rune(v)[0] == QueryNegationOperator {
		v = v[1:]
		isNegated = true
	}

	if len(v) > 0 && []rune(v)[len(v)-1] == QueryExactOperator {
		v = v[:len(v)-1]
		isExact = true
	}

	var p string

	if qp, ok := k.(kennung.QueryPrefixer); ok {
		p = qp.GetQueryPrefix()
	}

	if len(v) > 0 && v[:len(p)] == p {
		v = v[len(p):]
	}

	if err = k.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func KennungContainsExactlyMatchable(
	k kennung.KennungSansGattung,
	m Matchable,
) bool {
	switch kt := k.(type) {
	case kennung.EtikettLike:
		es := m.GetMetadatei().GetEtiketten()

		if es.ContainsKey(kt.GetEtikett().String()) {
			return true
		}

	case kennung.TypLike:
		if kennung.ContainsExactly(m.GetTyp(), k) {
			return true
		}

	default:
		// nop
	}

	idl := m.GetKennungLikePtr()

	if !kennung.ContainsExactly(idl, k) {
		return false
	}

	return true
}

func KennungContainsMatchable(
	k kennung.KennungSansGattung,
	m Matchable,
	ki kennung.Index,
) bool {
	switch kt := k.(type) {
	case kennung.EtikettLike:
		if iter.CheckAnyPtr[kennung.Etikett, *kennung.Etikett](
			m.GetMetadatei().GetEtiketten(),
			func(e *kennung.Etikett) (ok bool) {
				indexed, err := ki.Etiketten(e)

				var expanded kennung.EtikettSet

				if err == nil {
					expanded = indexed.GetExpandedRight()
				} else {
					expanded = kennung.ExpandOne(e, kennung.ExpanderRight)
				}

				ok = expanded.ContainsKey(expanded.KeyPtr(kt.GetEtikettPtr()))

				return
			},
		) {
			return true
		}

	case kennung.TypLike:
		if kennung.Contains(m.GetTyp(), k) {
			return true
		}

	case kennung.ShaLike:
		if Sha(kt.GetSha()).ContainsMatchable(m) {
			return true
		}

	case *kennung.Hinweis:
		// nop

	default:
		panic(fmt.Sprintf("unhandled type: %T", kt))
	}

	idl := m.GetKennungLikePtr()

	if !kennung.Contains(idl, k) {
		return false
	}

	return true
}