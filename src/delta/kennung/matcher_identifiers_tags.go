package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(Etikett) schnittstellen.Set[Etikett]
}

type MatcherIdentifierTags struct {
	Identifiers MatcherParentPtr
	Tags        MatcherParentPtr

	Matcher MatcherParentPtr
}

func MakeMatcherIdentifierTags() MatcherIdentifierTags {
	identifiers := MakeMatcherOrDoNotMatchOnEmpty()
	tags := MakeMatcherAndDoNotMatchOnEmpty()

	return MatcherIdentifierTags{
		Identifiers: identifiers,
		Tags:        tags,
		Matcher: MakeMatcherAnd(
			MakeMatcherOr(
				identifiers,
				tags,
			),
		),
	}
}

func tryAddMatcher(
	s *MatcherIdentifierTags,
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
			s.Identifiers.Add(m)
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

	err = errors.Wrap(errInvalidKennung(v))

	return
}

func (s MatcherIdentifierTags) MatcherLen() int {
	return s.Matcher.MatcherLen()
}

func (s MatcherIdentifierTags) String() string {
	return s.Matcher.String()
}

func (s *MatcherIdentifierTags) Add(m Matcher) (err error) {
	return s.Matcher.Add(m)
}

func (s MatcherIdentifierTags) ContainsMatchable(m Matchable) bool {
	return s.Matcher.ContainsMatchable(m)
}

func (s MatcherIdentifierTags) Len() int {
	return LenMatchers(s.Matcher) + s.Identifiers.MatcherLen()
}

func (s MatcherIdentifierTags) EachMatcher(
	f schnittstellen.FuncIter[Matcher],
) (err error) {
	return VisitAllMatchers(f, s.Matcher)
}

func (s MatcherIdentifierTags) GetHinweisen() schnittstellen.Set[Hinweis] {
	hins := collections.MakeMutableSetStringer[Hinweis]()

	VisitAllMatcherKennungSansGattungWrappers(
		func(m MatcherKennungSansGattungWrapper) (err error) {
			h, ok := m.GetKennung().(*Hinweis)

			if !ok {
				return
			}

			return hins.Add(*h)
		},
		s.Identifiers,
	)

	return hins
}

func (s MatcherIdentifierTags) AnyHinweis() (i1 Hinweis, ok bool) {
	if ok = s.Identifiers.MatcherLen() == 1; ok {
		hins := s.GetHinweisen()
		i1 = hins.Any()
	}

	return
}
