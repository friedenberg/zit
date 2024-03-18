package matcher_proto

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type QueryGroup interface {
	Get(g gattung.Gattung) (s MatcherSigil, ok bool)
	GetCwdFDs() fd.Set
	GetExplicitCwdFDs() fd.Set
	GetEtiketten() kennung.EtikettSet
	GetTypen() kennung.TypSet
	GetGattungen() gattungen.Set
	Matcher
}

type QueryGroupBuilder interface {
	BuildQueryGroup(...string) (QueryGroup, error)
}

type Matcher interface {
	ContainsMatchable(*sku.Transacted) bool
	schnittstellen.Stringer
	MatcherLen() int
	Each(schnittstellen.FuncIter[Matcher]) error
}

type MatcherSigil interface {
	Matcher
	GetSigil() kennung.Sigil
}

type MatcherSigilPtr interface {
	MatcherSigil
	AddSigil(kennung.Sigil)
}

type MatcherKennungSansGattungWrapper interface {
	Matcher
	GetKennung() kennung.KennungSansGattung
}

type MatcherExact interface {
	Matcher
	ContainsMatchableExactly(*sku.Transacted) bool
}

type MatcherImplicit interface {
	Matcher
	GetImplicitMatcher() Matcher
}

type MatcherParentPtr interface {
	Matcher
	Add(Matcher) error
}

type MatchableAdder interface {
	AddMatchable(*sku.Transacted) error
}

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(*kennung.Etikett) kennung.EtikettSet
}

type MatcherExactlyThisOrAllOfThese interface {
	Matcher
	AddExactlyThis(Matcher) error
	AddAllOfThese(Matcher) error
}

func LenMatchers(
	matchers ...Matcher,
) (i int) {
	inc := func(m Matcher) (err error) {
		if _, ok := m.(kennung.Kennung); ok {
			i++
		}

		return
	}

	VisitAllMatchers(inc, matchers...)

	return
}

func VisitAllMatcherKennungSansGattungWrappers(
	f schnittstellen.FuncIter[MatcherKennungSansGattungWrapper],
	ex func(Matcher) bool,
	matchers ...Matcher,
) (err error) {
	return VisitAllMatchers(
		func(m Matcher) (err error) {
			if ex != nil && ex(m) {
				return iter.MakeErrStopIteration()
			}

			if _, ok := m.(MatcherImplicit); ok {
				return iter.MakeErrStopIteration()
			}

			if mksgw, ok := m.(MatcherKennungSansGattungWrapper); ok {
				return f(mksgw)
			}

			return
		},
		matchers...,
	)
}

func VisitAllMatchers(
	f schnittstellen.FuncIter[Matcher],
	matchers ...Matcher,
) (err error) {
	for _, m := range matchers {
		if err = f(m); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if err = m.Each(
			func(m Matcher) (err error) {
				return VisitAllMatchers(f, m)
			},
		); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func SplitGattungenByHistory(qg QueryGroup) (schwanz, all kennung.Gattung) {
	err := qg.GetGattungen().Each(
		func(g gattung.Gattung) (err error) {
			m, ok := qg.Get(g)

			if !ok {
				return
			}

			if m.GetSigil().IncludesHistory() {
				all.Add(g)
			} else {
				schwanz.Add(g)
			}

			return
		},
	)

	errors.PanicIfError(err)

	return
}
