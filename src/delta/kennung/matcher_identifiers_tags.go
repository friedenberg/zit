package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(Etikett) schnittstellen.SetLike[Etikett]
}

type MatcherExactlyThisOrAllOfThese interface {
	Matcher
	AddExactlyThis(Matcher) error
	AddAllOfThese(Matcher) error
}

type matcherExactlyThisOrAllOfThese struct {
	MatcherExactlyThis MatcherParentPtr
	MatcherAllOfThese  MatcherParentPtr

	Matcher MatcherParentPtr
}

func MakeMatcherExactlyThisOrAllOfThese() MatcherExactlyThisOrAllOfThese {
	identifiers := MakeMatcherOrDoNotMatchOnEmpty()
	tags := MakeMatcherAndDoNotMatchOnEmpty()

	return &matcherExactlyThisOrAllOfThese{
		MatcherExactlyThis: identifiers,
		MatcherAllOfThese:  tags,
		Matcher: MakeMatcherOr(
			identifiers,
			tags,
		),
	}
}

func (s matcherExactlyThisOrAllOfThese) Each(
	f schnittstellen.FuncIter[Matcher],
) error {
	return s.Matcher.Each(f)
}

func (s *matcherExactlyThisOrAllOfThese) AddExactlyThis(m Matcher) (err error) {
	return s.MatcherExactlyThis.Add(m)
}

func (s *matcherExactlyThisOrAllOfThese) AddAllOfThese(m Matcher) (err error) {
	return s.MatcherAllOfThese.Add(m)
}

func (s matcherExactlyThisOrAllOfThese) MatcherLen() int {
	return s.Matcher.MatcherLen()
}

func (s matcherExactlyThisOrAllOfThese) String() string {
	return s.Matcher.String()
}

func (s *matcherExactlyThisOrAllOfThese) Add(m Matcher) (err error) {
	return s.Matcher.Add(m)
}

func (s matcherExactlyThisOrAllOfThese) ContainsMatchable(m Matchable) bool {
	return s.Matcher.ContainsMatchable(m)
}

func (s matcherExactlyThisOrAllOfThese) Len() int {
	return LenMatchers(s.Matcher) + s.MatcherExactlyThis.MatcherLen()
}

func (s matcherExactlyThisOrAllOfThese) EachMatcher(
	f schnittstellen.FuncIter[Matcher],
) (err error) {
	return VisitAllMatchers(f, s.Matcher)
}

func (s matcherExactlyThisOrAllOfThese) GetHinweisen() schnittstellen.SetLike[Hinweis] {
	hins := collections.MakeMutableSetStringer[Hinweis]()

	VisitAllMatcherKennungSansGattungWrappers(
		func(m MatcherKennungSansGattungWrapper) (err error) {
			h, ok := m.GetKennung().(*Hinweis)

			if !ok {
				return
			}

			return hins.Add(*h)
		},
		s.MatcherExactlyThis,
	)

	return hins
}

func (s matcherExactlyThisOrAllOfThese) AnyHinweis() (i1 Hinweis, ok bool) {
	if ok = s.MatcherExactlyThis.MatcherLen() == 1; ok {
		hins := s.GetHinweisen()
		i1 = hins.Any()
	}

	return
}
