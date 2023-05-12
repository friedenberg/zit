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
	expanders               Expanders
	implicitEtikettenGetter ImplicitEtikettenGetter

	Sigil Sigil

	UserMatcher   MatcherParentPtr
	ActualMatcher MatcherParentPtr
	cwd           matcherSigil
	hidden        matcherSigil
	Hinweisen     MatcherParentPtr
}

func MakeSet(
	cwd Matcher,
	ex Expanders,
	hidden Matcher,
	implicitEtikettenGetter ImplicitEtikettenGetter,
) Set {
	if hidden == nil {
		hidden = MakeMatcherNever()
	}

	return Set{
		expanders:               ex,
		implicitEtikettenGetter: implicitEtikettenGetter,
		UserMatcher:             MakeMatcherAnd(),
		ActualMatcher:           MakeMatcherAnd(),
		Hinweisen:               MakeMatcherOr(),
		cwd:                     MakeMatcherSigilMatchOnMissing(SigilCwd, cwd),
		hidden: MakeMatcherSigil(
			SigilHidden,
			MakeMatcherNegate(hidden),
		),
	}
}

func (s *Set) SetMany(vs ...string) (err error) {
	for _, v := range vs {
		if err = s.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Set) Set(v string) (err error) {
	{
		var m Matcher

		if m, err = MakeMatcher(&FD{}, v, nil); err == nil {
			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Sha{}, v, s.expanders.Sha); err == nil {
			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Konfig{}, v, nil); err == nil {
			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Hinweis{}, v, s.expanders.Hinweis); err == nil {
			s.Hinweisen.Add(m)
			return
		}
	}

	var (
		e         Etikett
		isNegated bool
	)

	if isNegated, err = SetQueryKennung(&e, v); err == nil {
		if s.implicitEtikettenGetter == nil {
			m := Matcher(e)

			if isNegated {
				m = MakeMatcherNegate(m)
			}

			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(m)
		} else {
			impl := s.implicitEtikettenGetter.GetImplicitEtiketten(e)

			mo := MakeMatcherOr()

			if err = impl.Each(
				func(e Etikett) (err error) {
					me := Matcher(e)

					if isNegated {
						me = MakeMatcherNegate(me)
					}

					return mo.Add(me)
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

			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(mo)
		}

		return
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Typ{}, v, s.expanders.Typ); err == nil {
			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(m)
			return
		}
	}

	{
		var m Matcher

		if m, err = MakeMatcher(&Kasten{}, v, s.expanders.Kasten); err == nil {
			s.UserMatcher.Add(m)
			s.ActualMatcher.Add(m)
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
			s.UserMatcher.Add(it)
			s.ActualMatcher.Add(it)

		case Sha:
			s.UserMatcher.Add(it)
			s.ActualMatcher.Add(it)

		case Hinweis:
			s.Hinweisen.Add(it)

		case Typ:
			s.UserMatcher.Add(it)
			s.ActualMatcher.Add(it)

		case Kasten:
			s.UserMatcher.Add(it)
			s.ActualMatcher.Add(it)

		case Konfig:
			s.UserMatcher.Add(it)
			s.ActualMatcher.Add(it)

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
	errors.TodoP4("improve the string creation method")
	sb := &strings.Builder{}

	errors.TodoP1("add Matchers")
	sb.WriteString(s.Sigil.String())

	return sb.String()
}

func (s Set) ContainsMatchable(m Matchable) bool {
	if !s.cwd.ContainsMatchable(m) {
		return false
	}

	if !s.hidden.ContainsMatchable(m) {
		return false
	}

	if !s.ActualMatcher.ContainsMatchable(m) {
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
	return LenMatchers(s.UserMatcher) + s.Hinweisen.Len()
}

func (s Set) EachMatcher(f schnittstellen.FuncIter[Matcher]) (err error) {
	return VisitAllMatchers(f, s.ActualMatcher)
}

func (s *Set) AddSigil(v Sigil) {
	s.Sigil.Add(v)
	s.hidden.Sigil.Add(v)
	s.cwd.Sigil.Add(v)
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
