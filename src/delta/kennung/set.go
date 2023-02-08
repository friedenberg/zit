package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/sha_collections"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Expanders struct {
	Sha, Etikett, Hinweis, Typ, Kasten func(string) (string, error)
}

type Set struct {
	expanders  Expanders
	Shas       sha_collections.MutableSet
	Etiketten  EtikettMutableSet
	Hinweisen  HinweisMutableSet
	Typen      TypMutableSet
	Timestamps ts.MutableSet
	Kisten     KastenMutableSet
	HasKonfig  bool
	Sigil      Sigil
}

func MakeSet(
	ex Expanders,
) Set {
	return Set{
		expanders:  ex,
		Shas:       sha_collections.MakeMutableSet(),
		Etiketten:  MakeEtikettMutableSet(),
		Hinweisen:  MakeHinweisMutableSet(),
		Typen:      MakeTypMutableSet(),
		Kisten:     MakeKastenMutableSet(),
		Timestamps: ts.MakeMutableSet(),
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

	if err = collections.ExpandAndAddString[Etikett, *Etikett](
		s.Etiketten,
		s.expanders.Etikett,
		v,
	); err == nil {
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

func (s *Set) Add(ids ...schnittstellen.Value) (err error) {
	for _, i := range ids {
		switch it := i.(type) {
		case Etikett:
			s.Etiketten.Add(it)

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

		default:
			err = errors.Errorf("unknown kennung: %s", it)
			return
		}
	}

	return
}

func (s Set) String() string {
	errors.TodoP4("improve the string creation method")
	return ""
}

func (s Set) OnlySingleHinweis() (h Hinweis, ok bool) {
	h = s.Hinweisen.Any()
	ok = s.Len() == 1 && s.Hinweisen.Len() == 1 && !s.Sigil.IncludesHistory()

	return
}

func (s Set) Len() int {
	k := 0

	if s.HasKonfig {
		k = 1
	}

	return s.Kisten.Len() + s.Shas.Len() + s.Etiketten.Len() + s.Hinweisen.Len() + s.Typen.Len() + s.Timestamps.Len() + k
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
