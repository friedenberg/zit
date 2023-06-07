package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(Etikett) schnittstellen.Set[Etikett]
}

type Set struct {
	Sigil Sigil

	Hinweisen MatcherParentPtr
	Others    MatcherParentPtr
	FDs       MatcherParentPtr

	Matcher MatcherParentPtr
}

func MakeSet() Set {
	hinweisen := MakeMatcherOrDoNotMatchOnEmpty()
	fds := MakeMatcherOrDoNotMatchOnEmpty()
	others := MakeMatcherAndDoNotMatchOnEmpty()

	return Set{
		Hinweisen: hinweisen,
		Others:    others,
		FDs:       fds,
		Matcher: MakeMatcherAnd(
			MakeMatcherOr(
				hinweisen,
				others,
				fds,
			),
		),
	}
}

func tryAddMatcher(
	s *Set,
	expanders Abbr,
	implicitEtikettenGetter ImplicitEtikettenGetter,
	v string,
) (err error) {
	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&FD{}, v, nil); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Sha{}, v, expanders.Sha.Expand); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Konfig{}, v, nil); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(
			&Hinweis{},
			v,
			expanders.Hinweis.Expand,
		); err == nil {
			s.Hinweisen.Add(m)
			return
		}
	}

	{
		var (
			e         Etikett
			isNegated bool
			// isExact   bool
			m Matcher
		)

		if m, isNegated, _, err = MakeMatcher(&e, v, nil); err == nil {
			if implicitEtikettenGetter == nil {
				s.Matcher.Add(m)
			} else {
				impl := implicitEtikettenGetter.GetImplicitEtiketten(e)

				mo := MakeMatcherOrDoNotMatchOnEmpty()

				if isNegated {
					mo = MakeMatcherAnd()
				}

				if err = impl.Each(
					func(e Etikett) (err error) {
						me := Matcher(MakeMatcherContainsExactly(e))

						if isNegated {
							me = MakeMatcherNegate(me)
						}

						return mo.Add(me)
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if isNegated {
					s.Matcher.Add(MakeMatcherAnd(m, MakeMatcherImplicit(mo)))
				} else {
					s.Matcher.Add(MakeMatcherOr(m, MakeMatcherImplicit(mo)))
				}
			}

			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Typ{}, v, expanders.Typ.Expand); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Kasten{}, v, expanders.Kasten.Expand); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	if err = s.Sigil.Set(v); err == nil {
		s.AddSigil(s.Sigil)
		return
	}

	err = errors.Wrap(errInvalidKennung(v))

	return
}

func (s Set) MatcherLen() int {
	return s.Matcher.MatcherLen()
}

func (s Set) String() string {
	return s.Matcher.String()
}

func (s *Set) Add(m Matcher) (err error) {
	return s.Matcher.Add(m)
}

func (s Set) ContainsMatchable(m Matchable) bool {
	return s.Matcher.ContainsMatchable(m)
}

func (s Set) Len() int {
	return LenMatchers(s.Matcher) + s.Hinweisen.MatcherLen()
}

func (s Set) EachMatcher(f schnittstellen.FuncIter[Matcher]) (err error) {
	return VisitAllMatchers(f, s.Matcher)
}

func (s *Set) AddSigil(v Sigil) {
	s.Sigil.Add(v)
}

func (s Set) GetSigil() Sigil {
	return s.Sigil
}

func (s Set) GetHinweisen() schnittstellen.Set[Hinweis] {
	hins := collections.MakeMutableSetStringer[Hinweis]()

	VisitAllMatcherKennungSansGattungWrappers(
		func(m MatcherKennungSansGattungWrapper) (err error) {
			h, ok := m.GetKennung().(*Hinweis)

			if !ok {
				return
			}

			return hins.Add(*h)
		},
		s.Hinweisen,
	)

	return hins
}

func (s Set) AnyHinweis() (i1 Hinweis, ok bool) {
	if ok = s.Hinweisen.MatcherLen() == 1; ok {
		hins := s.GetHinweisen()
		i1 = hins.Any()
	}

	return
}
