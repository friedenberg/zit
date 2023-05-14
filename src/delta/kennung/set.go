package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type ImplicitEtikettenGetter interface {
	GetImplicitEtiketten(Etikett) schnittstellen.Set[Etikett]
}

type Set struct {
	Sigil Sigil

	MatcherHidden MatcherSigilPtr
	MatcherCwd    MatcherSigilPtr
	Hinweisen     MatcherParentPtr
	Matcher       MatcherParentPtr
}

func MakeSet(
	cwd Matcher,
	hidden Matcher,
) Set {
	if hidden == nil {
		hidden = MakeMatcherNever()
	}

	sigilHidden := MakeMatcherSigil(
		SigilHidden,
		MakeMatcherNegate(hidden),
	)

	sigilCwd := MakeMatcherSigilMatchOnMissing(SigilCwd, cwd)

	return Set{
		Hinweisen:     MakeMatcherOr(),
		MatcherHidden: sigilHidden,
		MatcherCwd:    sigilCwd,
		Matcher: MakeMatcherImpExp(
			MakeMatcherAnd(sigilCwd, sigilHidden),
			MakeMatcherAnd(),
		),
	}
}

func tryAddMatcher(
	s *Set,
	expanders Expanders,
	implicitEtikettenGetter ImplicitEtikettenGetter,
	v string,
) (err error) {
	{
		var m Matcher

		if m, err = MakeMatcher(&FD{}, v, nil); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Sha{}, v, expanders.Sha); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Konfig{}, v, nil); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Hinweis{}, v, expanders.Hinweis); err == nil {
			s.Hinweisen.Add(m)
			return
		}
	}

	var (
		e         Etikett
		isNegated bool
	)

	if isNegated, err = SetQueryKennung(&e, v); err == nil {
		if implicitEtikettenGetter == nil {
			m := Matcher(e)

			if isNegated {
				m = MakeMatcherNegate(m)
			}

			s.Matcher.Add(m)
		} else {
			impl := implicitEtikettenGetter.GetImplicitEtiketten(e)

			mo := MakeMatcherOr()

			if isNegated {
				mo = MakeMatcherAnd()
			}

			if err = impl.Each(
				func(e Etikett) (err error) {
					me := Matcher(e)

					if isNegated {
						me = MakeMatcherNegate(me)
					}

					return mo.Add(MakeMatcherImplicit(me))
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			m := Matcher(e)

			if isNegated {
				m = MakeMatcherNegate(m)
			}

			mo.Add(m)

			s.Matcher.Add(mo)
		}

		return
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Typ{}, v, expanders.Typ); err == nil {
			s.Matcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Kasten{}, v, expanders.Kasten); err == nil {
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

func (s *Set) Add(ids ...schnittstellen.Element) (err error) {
	for _, i := range ids {
		switch it := i.(type) {
		case Etikett:
			errors.TodoP1("determine if this should have implicit etiketten")
			s.Matcher.Add(it)

		case Sha:
			s.Matcher.Add(it)

		case Hinweis:
			s.Matcher.Add(it)

		case Typ:
			s.Matcher.Add(it)

		case Kasten:
			s.Matcher.Add(it)

		case Konfig:
			s.Matcher.Add(it)

		case Sigil:
			s.AddSigil(it)

		case FD:
			if il, err := it.GetIdLike(); err == nil {
				s.Add(il)
			}

		default:
			err = errors.Errorf("unknown kennung: %s", it)
			return
		}
	}

	return
}

func (s Set) String() string {
	sb := &strings.Builder{}

	switch {
	case collections.Len(s.Matcher, s.Hinweisen) == 0:
		return s.Sigil.String()

	case s.Matcher.Len() == 0:
		sb.WriteString(QueryGroupOpenOperator)
		sb.WriteString(s.Hinweisen.String())
		sb.WriteString(QueryGroupCloseOperator)

	case s.Hinweisen.Len() == 0:
		sb.WriteString(QueryGroupOpenOperator)
		sb.WriteString(s.Matcher.String())
		sb.WriteString(QueryGroupCloseOperator)

	default:
		sb.WriteString(QueryGroupOpenOperator)
		sb.WriteString(s.Matcher.String())
		sb.WriteString(QueryOrOperator)
		sb.WriteString(s.Hinweisen.String())
		sb.WriteString(QueryGroupCloseOperator)

	}

	sb.WriteString(s.Sigil.String())

	return sb.String()
}

func (s Set) ContainsMatchable(m Matchable) bool {
	if !s.Matcher.ContainsMatchable(m) {
		return false
	}

	g := gattung.Must(m.GetGattung())

	if g != gattung.Zettel && s.Len() > 0 && s.Hinweisen.Len() == s.Len() {
		return false
	}

	il := m.GetIdLike()

	switch il.(type) {
	case Typ, Etikett, Kasten, Konfig:

	case Hinweis:
		if !s.Hinweisen.ContainsMatchable(m) {
			return false
		}

	default:
		panic(errors.Errorf("unsupported it type: %T, %s", il, il))
	}

	return true
}

func (s Set) Len() int {
	return LenMatchers(s.Matcher) + s.Hinweisen.Len()
}

func (s Set) EachMatcher(f schnittstellen.FuncIter[Matcher]) (err error) {
	return VisitAllMatchers(f, s.Matcher)
}

func (s *Set) AddSigil(v Sigil) {
	s.Sigil.Add(v)
	s.MatcherHidden.AddSigil(v)
	s.MatcherCwd.AddSigil(v)
}

func (s Set) GetSigil() schnittstellen.Sigil {
	return s.Sigil
}

func (s Set) GetHinweisen() schnittstellen.Set[Hinweis] {
	hins := collections.MakeMutableSetStringer[Hinweis]()

	VisitAllMatchers(
		func(m Matcher) (err error) {
			h, ok := m.(*Hinweis)

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
	if ok = s.Hinweisen.Len() == 1; ok {
		hins := s.GetHinweisen()
		i1 = hins.Any()
	}

	return
}
