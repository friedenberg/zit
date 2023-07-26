package kennung

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/gattungen"
)

func init() {
	gob.Register(&metaSet{})
}

// TODO-P3 rename to QueryGattungGroup
type MetaSet interface {
	Get(g gattung.Gattung) (s MatcherSigil, ok bool)
	GetFDs() schnittstellen.SetLike[FD]
	GetEtiketten() schnittstellen.SetLike[Etikett]
	GetTypen() schnittstellen.SetLike[Typ]
	Set(string) error
	SetMany(...string) error
	All(f func(gattung.Gattung, MatcherSigil) error) error
}

type setWithSigil struct {
	Matcher MatcherExactlyThisOrAllOfThese
	Sigil   Sigil
}

func (s setWithSigil) String() string {
	return fmt.Sprintf("%s%s", s.Matcher, s.Sigil)
}

func (s setWithSigil) ContainsMatchable(m Matchable) bool {
	return s.Matcher.ContainsMatchable(m)
}

func (s setWithSigil) GetSigil() Sigil {
	return s.Sigil
}

type metaSet struct {
	implicitEtikettenGetter ImplicitEtikettenGetter
	fileExtensionGetter     schnittstellen.FileExtensionGetter
	expanders               Abbr

	cwd    Matcher
	Hidden Matcher
	index  Index

	DefaultGattungen gattungen.Set
	Gattung          map[gattung.Gattung]setWithSigil
	FDs              schnittstellen.MutableSetLike[FD]
}

func MakeMetaSet(
	cwd Matcher,
	ex Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	dg gattungen.Set,
	implicitEtikettenGetter ImplicitEtikettenGetter,
	ki Index,
) MetaSet {
	return &metaSet{
		implicitEtikettenGetter: implicitEtikettenGetter,
		cwd:                     cwd,
		fileExtensionGetter:     feg,
		expanders:               ex,
		Hidden:                  hidden,
		DefaultGattungen:        dg.CloneMutableSetLike(),
		Gattung:                 make(map[gattung.Gattung]setWithSigil),
		FDs:                     collections.MakeMutableSetStringer[FD](),
		index:                   ki,
	}
}

func MakeMetaSetAll(
	cwd Matcher,
	ex Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	implicitEtikettenGetter ImplicitEtikettenGetter,
	ki Index,
) MetaSet {
	errors.TodoP2("support allowed sigils")
	return &metaSet{
		implicitEtikettenGetter: implicitEtikettenGetter,
		cwd:                     cwd,
		fileExtensionGetter:     feg,
		expanders:               ex,
		Hidden:                  hidden,
		DefaultGattungen:        gattungen.MakeSet(gattung.TrueGattung()...),
		Gattung:                 make(map[gattung.Gattung]setWithSigil),
		index:                   ki,
	}
}

func (s metaSet) String() string {
	sb := &strings.Builder{}

	hasAny := false

	for g, ids := range s.Gattung {
		if hasAny {
			sb.WriteString(" ")
		}

		hasAny = true

		sb.WriteString(
			fmt.Sprintf(
				"%s%s%s%s%s",
				QueryGroupOpenOperator,
				ids,
				QueryGroupCloseOperator,
				ids.Sigil,
				g,
			),
		)
	}

	return sb.String()
}

func (s *metaSet) SetMany(vs ...string) (err error) {
	builder := MatcherBuilder{}

	if _, err = builder.Build(vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range vs {
		if err = s.set(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	todo.Change("fix query syntax for groups")
	log.Log().Printf("query: %q", s)

	return
}

func (ms *metaSet) Set(v string) (err error) {
	return ms.set(v)
}

func (ms *metaSet) set(v string) (err error) {
	v = strings.TrimSpace(v)

	sbs := [3]*strings.Builder{
		{},
		{},
		{},
	}

	sbIdx := 0

	for _, c := range v {
		isSigil := SigilFieldFunc(c)

		switch {
		case isSigil && sbIdx == 0:
			sbIdx += 1

		case isSigil && sbIdx > 1:
			err = errors.Errorf("invalid meta set: %q", v)
			return

		case !isSigil && sbIdx == 1:
			sbIdx += 1
		}

		sbs[sbIdx].WriteRune(c)
	}

	var sigil Sigil

	if err = sigil.Set(sbs[1].String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	before := sbs[0].String()
	after := sbs[2].String()

	var gs gattungen.Set

	if after != "" {
		if gs, err = gattungen.GattungFromString(after); err != nil {
			if gattung.IsErrUnrecognizedGattung(err) {
				err = nil

				if err = collections.AddString[FD, *FD](
					ms.FDs,
					v,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

			} else {
				err = errors.Wrap(err)
			}

			return
		}
	} else {
		gs = ms.DefaultGattungen.CloneSetLike()
	}

	if err = gs.Each(
		func(g gattung.Gattung) (err error) {
			var ids setWithSigil
			ok := false

			if ids, ok = ms.Gattung[g]; !ok {
				ids.Matcher = MakeMatcherExactlyThisOrAllOfThese()
				ids.Sigil = sigil
			}

			switch {
			case before == "":
				break

			case ids.Sigil.IncludesCwd():
				fp := fmt.Sprintf("%s.%s", before, after)

				var fd FD

				if fd, err = FDFromPath(fp); err == nil {
					ids.Matcher.AddExactlyThis(fd)
					break
				}

				err = nil

				fallthrough

			default:
				if err = tryAddMatcher(
					ids.Matcher,
					ms.expanders,
					ms.implicitEtikettenGetter,
					ms.index,
					before,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			// if g.Equals(gattung.Konfig) {
			errors.TodoP1("move to gattung map")
			// ids.Matcher.Matcher.Add(MakeMatcherContainsExactly(Konfig{}))
			// }

			ms.Gattung[g] = ids

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func tryAddMatcher(
	s MatcherExactlyThisOrAllOfThese,
	expanders Abbr,
	implicitEtikettenGetter ImplicitEtikettenGetter,
	ki Index,
	v string,
) (err error) {
	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&FD{}, v, nil, ki); err == nil {
			return s.AddExactlyThis(m)
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Sha{}, v, expanders.Sha.Expand, ki); err == nil {
			return s.AddExactlyThis(m)
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(
			&Hinweis{},
			v,
			expanders.Hinweis.Expand,
			ki,
		); err == nil {
			s.AddExactlyThis(m)
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

		if m, isNegated, _, err = MakeMatcher(&e, v, nil, ki); err == nil {
			if implicitEtikettenGetter == nil {
				return s.AddAllOfThese(m)
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
					return s.AddAllOfThese(MakeMatcherAnd(m, MakeMatcherImplicit(mo)))
				} else {
					return s.AddAllOfThese(MakeMatcherOr(m, MakeMatcherImplicit(mo)))
				}
			}
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Typ{}, v, expanders.Typ.Expand, ki); err == nil {
			errors.TodoP1("handle typs that are incompatible")
			s.AddAllOfThese(m)
			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&Kasten{}, v, expanders.Kasten.Expand, ki); err == nil {
			return s.AddAllOfThese(m)
		}
	}

	err = errors.Wrap(errInvalidKennung(v))

	return
}

func (ms metaSet) Get(g gattung.Gattung) (s MatcherSigil, ok bool) {
	var ids setWithSigil

	ids, ok = ms.Gattung[g]

	sigilHidden := MakeMatcherExcludeHidden(ms.Hidden, ids.Sigil)
	sigilCwd := MakeMatcherSigilMatchOnMissing(SigilCwd, ms.cwd)

	s = MakeMatcherWithSigil(
		MakeMatcherAnd(
			MakeMatcherImplicit(sigilCwd),
			MakeMatcherImplicit(sigilHidden),
			ids.Matcher,
		),
		ids.Sigil,
	)

	return
}

func (ms metaSet) GetFDs() schnittstellen.SetLike[FD] {
	return ms.FDs
}

func (ms metaSet) GetEtiketten() schnittstellen.SetLike[Etikett] {
	es := MakeEtikettMutableSet()

	for _, s := range ms.Gattung {
		VisitAllMatcherKennungSansGattungWrappers(
			func(m MatcherKennungSansGattungWrapper) (err error) {
				e, ok := m.GetKennung().(EtikettLike)

				if !ok {
					return
				}

				return es.Add(e.GetEtikett())
			},
			// TODO-P1 modify sigil matcher to allow child traversal
			s.Matcher,
		)
	}

	return es
}

func (ms metaSet) GetTypen() schnittstellen.SetLike[Typ] {
	es := collections.MakeMutableSetStringer[Typ]()

	for _, s := range ms.Gattung {
		VisitAllMatcherKennungSansGattungWrappers(
			func(m MatcherKennungSansGattungWrapper) (err error) {
				e, ok := m.GetKennung().(Typ)

				if !ok {
					return
				}

				return es.Add(e)
			},
			// TODO-P1 modify sigil matcher to allow child traversal
			s.Matcher,
		)
	}

	return es
}

func (s metaSet) ContainsMatchable(m Matchable) bool {
	todo.Implement()
	return false
}

// Runs in parallel
func (ms metaSet) All(f func(gattung.Gattung, MatcherSigil) error) (err error) {
	errors.TodoP1("lock")
	chErr := make(chan error, len(ms.Gattung))

	for g := range ms.Gattung {
		ids, _ := ms.Get(g)

		go func(g gattung.Gattung, m MatcherSigil) {
			var err error

			defer func() {
				chErr <- err
			}()

			if err = f(g, ids); err != nil {
				if collections.IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}
		}(g, ids)
	}

	for i := 0; i < len(ms.Gattung); i++ {
		err = errors.Join(err, <-chErr)
	}

	return
}
