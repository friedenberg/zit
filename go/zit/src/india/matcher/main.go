package matcher

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/echo/kennung"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type MatchableAdder interface {
	AddMatchable(*sku.Transacted) error
}

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

	switch {
	case isNegated && isExact:
		m = MakeMatcherNegate(MakeMatcherContainsExactly(k))

	case isExact:
		m = MakeMatcherContainsExactly(k)

	case isNegated:
		m = MakeMatcherNegate(MakeMatcherContains(k, ki))

	default:
		m = MakeMatcherContains(k, ki)
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
	m *sku.Transacted,
) bool {
	switch kt := k.(type) {
	case kennung.Etikett, *kennung.Etikett:
		es := m.Metadatei.GetEtiketten()

		if es.ContainsKey(kt.String()) {
			return true
		}

	case kennung.TypLike:
		if kennung.ContainsExactly(m.GetTyp(), k) {
			return true
		}

	default:
		// nop
	}

	idl := &m.Kennung

	if !kennung.ContainsExactly(idl, k) {
		return false
	}

	return true
}

func KennungContainsMatchable(
	k kennung.KennungSansGattung,
	m *sku.Transacted,
	ki kennung.Index,
) bool {
	me := m.GetMetadatei()
	// log.Debug().Printf("%q -> %q", k, m.GetKennungLikePtr())

	switch kt := k.(type) {
	case kennung.Etikett, *kennung.Etikett:
		s := kt.String()

		if me.GetEtiketten().ContainsKey(s) {
			return true
		}

		if me.Verzeichnisse.GetExpandedEtiketten().ContainsKey(s) {
			return true
		}

		if me.Verzeichnisse.GetImplicitEtiketten().ContainsKey(s) {
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

	idl := &m.Kennung

	if !kennung.Contains(idl, k) {
		return false
	}

	return true
}
