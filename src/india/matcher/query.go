package matcher

import (
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func init() {
	gob.Register(&query{})
}

// TODO-P3 rename to QueryGattungGroup
type Query interface {
	Get(g gattung.Gattung) (s MatcherSigil, ok bool)
	GetCwdFDs() fd.Set
	GetExplicitCwdFDs() fd.Set
	GetEtiketten() kennung.EtikettSet
	GetTypen() schnittstellen.SetLike[kennung.Typ]
	Set(string) error
	SetMany(...string) error
	All(f func(gattung.Gattung, MatcherSigil) error) error
	GetGattungen() gattungen.Set
	Matcher
}

func MakeQueryFromCheckedOutSet(
	cos sku.CheckedOutSet,
) (q Query, err error) {
	gs := make(map[gattung.Gattung]setWithSigil)

	if err = cos.EachPtr(
		func(co *sku.CheckedOut) (err error) {
			m := MakeMatcherContainsExactly(co.Internal.Kennung)

			var s setWithSigil
			ok := false

			if s, ok = gs[gattung.Must(co.Internal.GetGattung())]; !ok {
				s.Matcher = MakeMatcherExactlyThisOrAllOfThese()
			}

			s.Matcher.AddExactlyThis(m)

			gs[gattung.Must(co.Internal.GetGattung())] = s

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	q = &query{
		Gattung: gs,
	}

	return
}

func MakeQuery(
	k schnittstellen.Konfig,
	cwd MatcherCwd,
	ex kennung.Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	dg gattungen.Set,
	ki kennung.Index,
) Query {
	return &query{
		konfig:              k,
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Hidden:              hidden,
		DefaultGattungen:    dg.CloneMutableSetLike(),
		Gattung:             make(map[gattung.Gattung]setWithSigil),
		FDs:                 fd.MakeMutableSet(),
		index:               ki,
	}
}

func MakeQueryAll(
	k schnittstellen.Konfig,
	cwd MatcherCwd,
	ex kennung.Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	ki kennung.Index,
) Query {
	errors.TodoP2("support allowed sigils")
	return &query{
		konfig:              k,
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Hidden:              hidden,
		DefaultGattungen:    gattungen.MakeSet(gattung.TrueGattung()...),
		Gattung:             make(map[gattung.Gattung]setWithSigil),
		index:               ki,
	}
}

type setWithSigil struct {
	Matcher MatcherExactlyThisOrAllOfThese
	Sigil   kennung.Sigil
}

func (s setWithSigil) String() string {
	return fmt.Sprintf("%s%s", s.Matcher, s.Sigil)
}

func (s setWithSigil) ContainsMatchable(m *sku.Transacted) bool {
	return s.Matcher.ContainsMatchable(m)
}

func (s setWithSigil) GetSigil() kennung.Sigil {
	return s.Sigil
}

type query struct {
	konfig              schnittstellen.Konfig
	fileExtensionGetter schnittstellen.FileExtensionGetter
	expanders           kennung.Abbr

	cwd    MatcherCwd
	Hidden Matcher
	index  kennung.Index

	DefaultGattungen gattungen.Set
	Gattung          map[gattung.Gattung]setWithSigil
	FDs              fd.MutableSet

	dotOperatorActive bool
}

func (q query) MatcherLen() int {
	return 0
}

func (q query) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	for _, s := range q.Gattung {
		if err = f(s.Matcher); err != nil {
			return
		}
	}

	return
}

func (s query) String() string {
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

func (s *query) SetMany(vs ...string) (err error) {
	builder := MatcherBuilder{}

	if _, err = builder.Build(vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range vs {
		if err = s.set(v); err != nil {
			var fd fd.FD

			if err1 := fd.Set(v); err1 == nil {
				if err = s.FDs.Add(fd); err != nil {
					err = errors.Wrap(err)
					return
				}

				continue
			}

			err = errors.Wrap(err)
			return
		}
	}

	todo.Change("fix query syntax for groups")
	log.Log().Printf("query: %q", s)

	return
}

func (ms *query) Set(v string) (err error) {
	return ms.set(v)
}

func (ms *query) set(v string) (err error) {
	v = strings.TrimSpace(v)

	sbs := [3]*strings.Builder{
		{},
		{},
		{},
	}

	sbIdx := 0

	for _, c := range v {
		isSigil := kennung.SigilFieldFunc(c)

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

	var sigil kennung.Sigil

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

				if err = iter.AddString[fd.FD, *fd.FD](
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

	if before == "" && after == "" && sigil.IncludesCwd() {
		ms.dotOperatorActive = true
		// return
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

				var f fd.FD

				if f, err = fd.FDFromPath(fp); err == nil {
					ids.Matcher.AddExactlyThis(FD(f))
					break
				}

				err = nil

				fallthrough

			default:
				if err = tryAddMatcher(
					ms.konfig,
					ids.Matcher,
					ms.expanders,
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
	k schnittstellen.Konfig,
	s MatcherExactlyThisOrAllOfThese,
	expanders kennung.Abbr,
	ki kennung.Index,
	v string,
) (err error) {
	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&fd.FD{}, v, nil, ki, k); err == nil {
			return s.AddExactlyThis(m)
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&kennung.Sha{}, v, expanders.Sha.Expand, ki, k); err == nil {
			return s.AddExactlyThis(m)
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(
			&kennung.Hinweis{},
			v,
			expanders.Hinweis.Expand,
			ki,
			k,
		); err == nil {
			s.AddExactlyThis(m)
			return
		}
	}

	if lString, ok := k.GetFilters()[v]; ok {
		var m Matcher

		if m, err = MakeMatcherLua(ki, lString); err == nil {
			s.AddAllOfThese(m)
			return
		}
	}

	{
		var (
			e         kennung.Etikett
			isNegated bool
			// isExact   bool
			m Matcher
		)

		if m, isNegated, _, err = MakeMatcher(&e, v, nil, ki, k); err == nil {
			mo := MakeMatcherOrDoNotMatchOnEmpty()

			if isNegated {
				mo = MakeMatcherAnd()
			}

			if isNegated {
				return s.AddAllOfThese(MakeMatcherAnd(m, MakeMatcherImplicit(mo)))
			} else {
				return s.AddAllOfThese(MakeMatcherOr(m, MakeMatcherImplicit(mo)))
			}
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&kennung.Typ{}, v, expanders.Typ.Expand, ki, k); err == nil {
			errors.TodoP1("handle typs that are incompatible")
			s.AddAllOfThese(m)
			return
		}
	}

	{
		var m Matcher

		if m, _, _, err = MakeMatcher(&kennung.Kasten{}, v, expanders.Kasten.Expand, ki, k); err == nil {
			return s.AddAllOfThese(m)
		}
	}

	err = errors.Wrap(kennung.ErrInvalidKennung(v))

	return
}

func (ms query) Get(g gattung.Gattung) (s MatcherSigil, ok bool) {
	var ids setWithSigil

	ids, ok = ms.Gattung[g]

	sigilHidden := MakeMatcherExcludeHidden(ms.Hidden, ids.Sigil)
	sigilCwd := MakeMatcherSigilMatchOnMissing(kennung.SigilCwd, ms.cwd)

	s = MakeMatcherWithSigil(
		MakeMatcherAnd(
			ids.Matcher,
			MakeMatcherImplicit(sigilCwd),
			MakeMatcherImplicit(sigilHidden),
		),
		ids.Sigil,
	)

	return
}

func (ms query) GetExplicitCwdFDs() fd.Set {
	return ms.FDs
}

func (ms query) GetCwdFDs() fd.Set {
	if ms.dotOperatorActive {
		return ms.cwd.GetCwdFDs()
	} else {
		return ms.FDs
	}
}

func (ms query) GetEtiketten() kennung.EtikettSet {
	es := kennung.MakeEtikettMutableSet()

	for _, s := range ms.Gattung {
		VisitAllMatcherKennungSansGattungWrappers(
			func(m MatcherKennungSansGattungWrapper) (err error) {
				e, ok := m.GetKennung().(kennung.EtikettLike)

				if !ok {
					return
				}

				return es.AddPtr(e.GetEtikettPtr())
			},
			IsMatcherNegate,
			// TODO-P1 modify sigil matcher to allow child traversal
			s.Matcher,
		)
	}

	return es
}

func (ms query) GetTyp() (t kennung.Typ, ok bool) {
	ts := ms.GetTypen()

	if ts.Len() != 1 {
		return
	}

	t = ts.Any()
	ok = true

	return
}

func (ms query) GetTypen() schnittstellen.SetLike[kennung.Typ] {
	es := kennung.MakeMutableTypSet()

	for _, s := range ms.Gattung {
		VisitAllMatcherKennungSansGattungWrappers(
			func(m MatcherKennungSansGattungWrapper) (err error) {
				e, ok := m.GetKennung().(kennung.TypLike)

				if !ok {
					return
				}

				return es.AddPtr(e.GetTypPtr())
			},
			func(m Matcher) bool {
				ok := false

				switch m.(type) {
				case Negate, *Negate:
					ok = true
				}

				return ok
			},
			// TODO-P1 modify sigil matcher to allow child traversal
			s.Matcher,
		)
	}

	return es
}

func (s query) ContainsMatchable(m *sku.Transacted) bool {
	g := gattung.Must(m.GetGattung())

	var matcher Matcher
	ok := false

	if matcher, ok = s.Get(g); !ok {
		return false
	}

	return matcher.ContainsMatchable(m)
}

func (ms query) GetGattungen() gattungen.Set {
	gs := make([]gattung.Gattung, 0, len(ms.Gattung))

	for g := range ms.Gattung {
		gs = append(gs, g)
	}

	return gattungen.MakeSet(gs...)
}

// Runs in parallel
func (ms query) All(f func(gattung.Gattung, MatcherSigil) error) (err error) {
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
				if iter.IsStopIteration(err) {
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
