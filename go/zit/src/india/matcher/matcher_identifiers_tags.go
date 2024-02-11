package matcher

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(*kennung.Etikett) kennung.EtikettSet
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

func (s matcherExactlyThisOrAllOfThese) ContainsMatchable(m *sku.Transacted) bool {
	ok := s.Matcher.ContainsMatchable(m)
	// log.Debug().Printf("%q -> %t -> %q", s.Matcher, ok, m)
	return ok
}

func (s matcherExactlyThisOrAllOfThese) Len() int {
	return LenMatchers(s.Matcher) + s.MatcherExactlyThis.MatcherLen()
}

func (s matcherExactlyThisOrAllOfThese) EachMatcher(
	f schnittstellen.FuncIter[Matcher],
) (err error) {
	return VisitAllMatchers(f, s.Matcher)
}

func (s matcherExactlyThisOrAllOfThese) GetHinweisen() schnittstellen.SetLike[kennung.Hinweis] {
	hins := collections_value.MakeMutableValueSet[kennung.Hinweis](nil)

	VisitAllMatcherKennungSansGattungWrappers(
		func(m MatcherKennungSansGattungWrapper) (err error) {
			h, ok := m.GetKennung().(*kennung.Hinweis)

			if !ok {
				return
			}

			return hins.Add(*h)
		},
		func(m Matcher) bool {
			ok := false

			switch m.(type) {
			case Negate, *Negate:
				ok = true
			}

			return ok
		},
		s.MatcherExactlyThis,
	)

	return hins
}

func (s matcherExactlyThisOrAllOfThese) AnyHinweis() (i1 kennung.Hinweis, ok bool) {
	if ok = s.MatcherExactlyThis.MatcherLen() == 1; ok {
		hins := s.GetHinweisen()
		i1 = hins.Any()
	}

	return
}
