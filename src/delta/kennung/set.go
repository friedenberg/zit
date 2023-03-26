package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/sha_collections"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Set struct {
	cwd        Matcher
	expanders  Expanders
	Shas       sha_collections.MutableSet
	Etiketten  MutableQuerySet[Etikett, *Etikett]
	Hinweisen  HinweisMutableSet
	Typen      TypMutableSet
	Timestamps ts.MutableSet
	Kisten     KastenMutableSet
	FDs        MutableFDSet
	HasKonfig  bool
	Sigil      Sigil

	hidden Matcher
}

func MakeSet(
	cwd Matcher,
	ex Expanders,
	hidden Matcher,
) Set {
	if hidden == nil {
		hidden = MakeMatcherNever()
	}

	return Set{
		cwd:        cwd,
		expanders:  ex,
		Shas:       sha_collections.MakeMutableSet(),
		Etiketten:  MakeMutableQuerySet[Etikett, *Etikett](ex.Etikett, nil, nil),
		Hinweisen:  MakeHinweisMutableSet(),
		Typen:      MakeTypMutableSet(),
		Kisten:     MakeKastenMutableSet(),
		Timestamps: ts.MakeMutableSet(),
		FDs:        MakeMutableFDSet(),
		hidden:     hidden,
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
	if err = collections.AddString[FD, *FD](
		s.FDs,
		v,
	); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[sha.Sha, *sha.Sha](
		s.Shas,
		s.expanders.Sha,
		v,
	); err == nil {
		return
	}

	if err = collections.AddString[ts.Time, *ts.Time](s.Timestamps, v); err == nil {
		return
	}

	if err = (Konfig{}).Set(v); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[Hinweis, *Hinweis](
		s.Hinweisen,
		s.expanders.Hinweis,
		v,
	); err == nil {
		return
	}

	if err = s.Etiketten.AddString(v); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[Typ, *Typ](
		s.Typen,
		s.expanders.Typ,
		v,
	); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[Kasten, *Kasten](
		s.Kisten,
		s.expanders.Kasten,
		v,
	); err == nil {
		return
	}

	if err = s.Sigil.Set(v); err == nil {
		return
	}

	err = errors.Wrap(errInvalidKennung(v))

	return
}

func (s *Set) Add(ids ...schnittstellen.Element) (err error) {
	for _, i := range ids {
		switch it := i.(type) {
		case Etikett:
			s.Etiketten.AddInclude(it)

		case sha.Sha:
			s.Shas.Add(it)

		case Hinweis:
			s.Hinweisen.Add(it)

		case Typ:
			s.Typen.Add(it)

		case ts.Time:
			s.Timestamps.Add(it)

		case Kasten:
			s.Kisten.Add(it)

		case Konfig:
			s.HasKonfig = true

		case Sigil:
			s.Sigil.Add(it)

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

	s.Shas.Each(iter.AddString[sha.Sha](sb))
	sb.WriteString(s.Etiketten.String())
	sb.WriteString(" ")
	s.Hinweisen.Each(iter.AddString[Hinweis](sb))
	s.Typen.Each(iter.AddString[Typ](sb))
	s.Timestamps.Each(iter.AddString[ts.Time](sb))
	s.Kisten.Each(iter.AddString[Kasten](sb))
	s.FDs.Each(iter.AddString[FD](sb))

	if s.HasKonfig {
		sb.WriteString("konfig")
	}

	sb.WriteString(s.Sigil.String())

	return sb.String()
}

func (s Set) ContainsMatchable(m Matchable) bool {
	if s.cwd != nil && !s.cwd.ContainsMatchable(m) {
		return false
	}

	if !s.Sigil.IncludesHidden() &&
		s.hidden != nil &&
		s.hidden.ContainsMatchable(m) {
		return false
	}

	g := gattung.Must(m.GetGattung())

	es := m.GetEtikettenExpanded()
	containsEtts := s.Etiketten.ContainsAgainst(es)

	if !containsEtts && g != gattung.Etikett {
		return false
	}

	// Only Zettels have Typs, so only filter against them in that case
	if g == gattung.Zettel {
		if t := m.GetTyp(); s.Typen.Len() > 0 && !s.Typen.Contains(t) {
			return false
		}
		// If this is a strict Hinweis match, do not permit this to match anything
		// other than Zettels
	} else if s.Len() > 0 && s.Hinweisen.Len() == s.Len() {
		return false
	}

	il := m.GetIdLike()

	switch id := il.(type) {
	case Kasten:

	case Typ:
		if s.Typen.Len() > 0 && !s.Typen.Contains(id) {
			return false
		}

	case Etikett:
		if s.Etiketten.Len() > 0 && !s.Etiketten.Contains(id) {
			return false
		}

	case Hinweis:
		if s.Hinweisen.Len() > 0 && !s.Hinweisen.Contains(id) {
			return false
		}

	default:
		panic(errors.Errorf("unsupported it type: %T, %s", il, il))
	}

	return true
}

// func (s Set) Contains(id schnittstellen.Stringer) bool {
// 	switch idt := id.(type) {
// 	case sha.Sha:
// 		return s.Shas.Contains(idt)

// 	case Etikett:
// 		return s.Etiketten.Contains(idt)

// 	case Typ:
// 		return s.Typen.Contains(idt)

// 	case *Hinweis:
// 		return s.Hinweisen.Contains(*idt)

// 	case Hinweis:
// 		return s.Hinweisen.Contains(idt)

// 	case ts.Time:
// 		return s.Timestamps.Contains(idt)

// 	case Kasten:
// 		return s.Kisten.Contains(idt)

// 	case FD:
// 		return s.FDs.Contains(idt)

// 	case Konfig:
// 		return true

// 	default:
// 		return false
// 	}
// }

func (s Set) OnlySingleHinweis() (h Hinweis, ok bool) {
	if s.Len() != 1 {
		return
	}

	if s.Sigil.IncludesHistory() {
		return
	}

	switch {
	case s.Hinweisen.Len() == 1:
		ok = true
		h = s.Hinweisen.Any()

	case s.FDs.Len() == 1:
		var err error

		h, err = s.FDs.Any().GetHinweis()
		ok = err == nil

	default:
	}

	return
}

func (s Set) Len() int {
	k := 0

	if s.HasKonfig {
		k = 1
	}

	return collections.Len(
		s.Kisten,
		s.Shas,
		s.Etiketten,
		s.Hinweisen,
		s.Typen,
		s.Timestamps,
		s.FDs,
	) + k
}

func (s Set) GetSigil() schnittstellen.Sigil {
	return s.Sigil
}

func (s Set) AnyShasOrHinweisen() (ids []schnittstellen.Korper) {
	ids = make([]schnittstellen.Korper, 0, s.Shas.Len()+s.Hinweisen.Len())

	s.Shas.Each(
		func(sh sha.Sha) (err error) {
			ids = append(ids, sh)

			return
		},
	)

	s.Hinweisen.Each(
		func(h Hinweis) (err error) {
			ids = append(ids, h)

			return
		},
	)

	return
}

func (s Set) AnyShaOrHinweis() (i1 schnittstellen.Korper, ok bool) {
	ids := s.AnyShasOrHinweisen()

	if len(ids) > 0 {
		i1 = ids[0]
		ok = true
	}

	return
}

func (s Set) AnyHinweis() (i1 Hinweis, ok bool) {
	if ok = s.Hinweisen.Len() == 1; ok {
		i1 = s.Hinweisen.Any()
	}

	return
}
