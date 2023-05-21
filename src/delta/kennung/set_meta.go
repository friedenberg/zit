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
	Get(g gattung.Gattung) (s Matcher, ok bool)
	GetFDs() schnittstellen.Set[FD]
	GetEtiketten() schnittstellen.Set[Etikett]
	GetTypen() schnittstellen.Set[Typ]
	Set(string) error
	SetMany(...string) error
	All(f func(gattung.Gattung, Matcher) error) error
}

type metaSet struct {
	implicitEtikettenGetter ImplicitEtikettenGetter
	fileExtensionGetter     schnittstellen.FileExtensionGetter
	expanders               Abbr

	cwd    Matcher
	Hidden Matcher

	DefaultGattungen gattungen.Set
	Gattung          map[gattung.Gattung]Set
	FDs              schnittstellen.MutableSet[FD]
}

func MakeMetaSet(
	cwd Matcher,
	ex Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	dg gattungen.Set,
	implicitEtikettenGetter ImplicitEtikettenGetter,
) MetaSet {
	return &metaSet{
		implicitEtikettenGetter: implicitEtikettenGetter,
		cwd:                     cwd,
		fileExtensionGetter:     feg,
		expanders:               ex,
		Hidden:                  hidden,
		DefaultGattungen:        dg.MutableClone(),
		Gattung:                 make(map[gattung.Gattung]Set),
		FDs:                     collections.MakeMutableSetStringer[FD](),
	}
}

func MakeMetaSetAll(
	cwd Matcher,
	ex Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	implicitEtikettenGetter ImplicitEtikettenGetter,
) MetaSet {
	errors.TodoP2("support allowed sigils")
	return &metaSet{
		implicitEtikettenGetter: implicitEtikettenGetter,
		cwd:                     cwd,
		fileExtensionGetter:     feg,
		expanders:               ex,
		Hidden:                  hidden,
		DefaultGattungen:        gattungen.MakeSet(gattung.TrueGattung()...),
		Gattung:                 make(map[gattung.Gattung]Set),
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
		sb.WriteString(fmt.Sprintf("%s%s", ids, g))
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
		gs = ms.DefaultGattungen.ImmutableClone()
	}

	if err = gs.Each(
		func(g gattung.Gattung) (err error) {
			var ids Set
			ok := false

			if ids, ok = ms.Gattung[g]; !ok {
				ids = ms.MakeSet()
				ids.AddSigil(sigil)
			}

			switch {
			case before == "":
				break

			case ids.Sigil.IncludesCwd():
				fp := fmt.Sprintf("%s.%s", before, after)

				var fd FD

				if fd, err = FDFromPath(fp); err == nil {
					ids.Add(fd)
					break
				}

				err = nil

				fallthrough

			default:
				if err = tryAddMatcher(
					&ids,
					ms.expanders,
					ms.implicitEtikettenGetter,
					before,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if g.Equals(gattung.Konfig) {
				ids.Add(Konfig{})
			}

			ms.Gattung[g] = ids

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ms metaSet) Get(g gattung.Gattung) (s Matcher, ok bool) {
	s, ok = ms.Gattung[g]

	hidden := ms.Hidden

	if hidden == nil {
		hidden = MakeMatcherNever()
	}

	sigilHidden := MakeMatcherSigil(
		SigilHidden,
		MakeMatcherNegate(hidden),
	)

	sigilCwd := MakeMatcherSigilMatchOnMissing(SigilCwd, ms.cwd)

	s = MakeMatcherImpExp(
		MakeMatcherAnd(sigilCwd, sigilHidden),
		MakeMatcherAnd(s),
	)

	return
}

func (ms metaSet) GetFDs() schnittstellen.Set[FD] {
	return ms.FDs
}

func (ms metaSet) GetEtiketten() schnittstellen.Set[Etikett] {
	es := MakeEtikettMutableSet()

	for _, s := range ms.Gattung {
		VisitAllMatchers(
			func(m Matcher) (err error) {
				e, ok := m.(Etikett)

				if !ok {
					return
				}

				return es.Add(e)
			},
			s.Matcher,
		)
	}

	return es
}

func (ms metaSet) GetTypen() schnittstellen.Set[Typ] {
	es := collections.MakeMutableSetStringer[Typ]()

	for _, s := range ms.Gattung {
		VisitAllMatchers(
			func(m Matcher) (err error) {
				e, ok := m.(Typ)

				if !ok {
					return
				}

				return es.Add(e)
			},
			s.Matcher,
		)
	}

	return es
}

func (ms metaSet) MakeSet() Set {
	return MakeSet(ms.cwd, ms.Hidden)
}

func (s metaSet) ContainsMatchable(m Matchable) bool {
	todo.Implement()
	return false
}

// Runs in parallel
func (ms metaSet) All(f func(gattung.Gattung, Matcher) error) (err error) {
	errors.TodoP1("lock")
	chErr := make(chan error, len(ms.Gattung))

	for g, s := range ms.Gattung {
		go func(g gattung.Gattung, ids Set) {
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
		}(g, s)
	}

	for i := 0; i < len(ms.Gattung); i++ {
		err = errors.Join(err, <-chErr)
	}

	return
}
