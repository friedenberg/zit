package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(Etikett) schnittstellen.Set[Etikett]
}

type MatcherIdentifierTags struct {
	MatcherExactlyThis MatcherParentPtr
	MatcherAllOfThese  MatcherParentPtr

	Matcher MatcherParentPtr
}

func MakeMatcherIdentifierTags() MatcherIdentifierTags {
	identifiers := MakeMatcherOrDoNotMatchOnEmpty()
	tags := MakeMatcherAndDoNotMatchOnEmpty()

	return MatcherIdentifierTags{
		MatcherExactlyThis: identifiers,
		MatcherAllOfThese:  tags,
		Matcher: MakeMatcherOr(
			identifiers,
			tags,
		),
	}
}

func (s *MatcherIdentifierTags) AddExactlyThis(m Matcher) (err error) {
	return s.MatcherExactlyThis.Add(m)
}

func (s *MatcherIdentifierTags) AddAllOfThese(m Matcher) (err error) {
	return s.MatcherAllOfThese.Add(m)
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
	return LenMatchers(s.Matcher) + s.MatcherExactlyThis.MatcherLen()
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
		s.MatcherExactlyThis,
	)

	return hins
}

func (s MatcherIdentifierTags) AnyHinweis() (i1 Hinweis, ok bool) {
	if ok = s.MatcherExactlyThis.MatcherLen() == 1; ok {
		hins := s.GetHinweisen()
		i1 = hins.Any()
	}

	return
}
