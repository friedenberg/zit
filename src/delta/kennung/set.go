package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/sha_collections"
	"github.com/friedenberg/zit/src/echo/ts"
)

// TODO-P3 rewrite
type Set struct {
	shaExpander func(string) (string, error)
	Shas        sha_collections.MutableSet

	etikettExpander func(string) (string, error)
	Etiketten       EtikettMutableSet

	hinweisExpander func(string) (string, error)
	Hinweisen       HinweisMutableSet

	typExpander func(string) (string, error)
	Typen       TypMutableSet

	Timestamps ts.MutableSet

	kastenExpander func(string) (string, error)
	Kisten         KastenMutableSet

	HasKonfig bool
	Sigil     Sigil
}

func MakeSet() Set {
	return Set{
		Timestamps: ts.MakeMutableSet(),
		Shas:       sha_collections.MakeMutableSet(),
		Etiketten:  MakeEtikettMutableSet(),
		Hinweisen:  MakeHinweisMutableSet(),
		Typen:      MakeTypMutableSet(),
		Kisten:     MakeKastenMutableSet(),
	}
}

func MakeSetWithExpanders(
	shaExpander func(string) (string, error),
	etikettExpander func(string) (string, error),
	hinweisExpander func(string) (string, error),
	typExpander func(string) (string, error),
	kastenExpander func(string) (string, error),
) Set {
	errors.TodoP0("implement")
	return Set{
		shaExpander:     shaExpander,
		Shas:            sha_collections.MakeMutableSet(),
		etikettExpander: etikettExpander,
		Etiketten:       MakeEtikettMutableSet(),
		hinweisExpander: hinweisExpander,
		Hinweisen:       MakeHinweisMutableSet(),
		typExpander:     typExpander,
		Typen:           MakeTypMutableSet(),
		kastenExpander:  kastenExpander,
		Kisten:          MakeKastenMutableSet(),
		Timestamps:      ts.MakeMutableSet(),
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
		s.shaExpander,
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
		s.hinweisExpander,
		v,
	); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[Etikett, *Etikett](
		s.Etiketten,
		s.etikettExpander,
		v,
	); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[Typ, *Typ](
		s.Typen,
		s.typExpander,
		v,
	); err == nil {
		return
	}

	if err = collections.ExpandAndAddString[Kasten, *Kasten](
		s.Kisten,
		s.kastenExpander,
		v,
	); err == nil {
		return
	}

	if err = s.Sigil.Set(v); err == nil {
		return
	}

	err = errors.Wrap(ErrInvalid)

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
	ok = s.Len() == 1 && s.Hinweisen.Len() == 1 && !h.GetSigil().IncludesHistory()

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
