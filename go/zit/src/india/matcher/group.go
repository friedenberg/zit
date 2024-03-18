package matcher

import (
	"encoding/gob"
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/gattungen"
	"code.linenisgreat.com/zit/src/delta/zittish"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
)

func init() {
	gob.Register(&group{})
}

type Group interface {
	matcher_proto.QueryGroup
}

func MakeGroupFromCheckedOutSet(
	cos sku.CheckedOutSet,
) (q Group, err error) {
	gs := make(map[gattung.Gattung]Query)

	if err = cos.Each(
		func(co *sku.CheckedOut) (err error) {
			m := MakeMatcherContainsExactly(&co.Internal.Kennung)

			var s Query
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

	q = &group{
		Gattung: gs,
	}

	return
}

func MakeGroup(
	k schnittstellen.Konfig,
	cwd matcher_proto.Cwd,
	ex kennung.Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	dg kennung.Gattung,
	ki kennung.Index,
) matcher_proto.QueryGroupBuilder {
	return &group{
		konfig:              k,
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Hidden:              hidden,
		DefaultGattungen:    dg,
		Gattung:             make(map[gattung.Gattung]Query),
		FDs:                 fd.MakeMutableSet(),
		index:               ki,
	}
}

func MakeGroupAll(
	k schnittstellen.Konfig,
	cwd matcher_proto.Cwd,
	ex kennung.Abbr,
	hidden Matcher,
	feg schnittstellen.FileExtensionGetter,
	ki kennung.Index,
) Group {
	errors.TodoP2("support allowed sigils")
	return &group{
		konfig:              k,
		cwd:                 cwd,
		fileExtensionGetter: feg,
		expanders:           ex,
		Hidden:              hidden,
		DefaultGattungen:    kennung.MakeGattung(gattung.TrueGattung()...),
		Gattung:             make(map[gattung.Gattung]Query),
		index:               ki,
	}
}

type group struct {
	konfig              schnittstellen.Konfig
	fileExtensionGetter schnittstellen.FileExtensionGetter
	expanders           kennung.Abbr

	cwd    matcher_proto.Cwd
	Hidden Matcher
	index  kennung.Index

	DefaultGattungen kennung.Gattung
	Gattung          map[gattung.Gattung]Query
	// NewQuery         *query.QueryGroup
	FDs fd.MutableSet

	dotOperatorActive bool
}

func (q group) MatcherLen() int {
	return 0
}

func (q group) Each(f schnittstellen.FuncIter[Matcher]) (err error) {
	for _, s := range q.Gattung {
		if err = f(s.Matcher); err != nil {
			return
		}
	}

	return
}

func (s group) String() string {
	sb := &strings.Builder{}

	hasAny := false

	for g, ids := range s.Gattung {
		if hasAny {
			sb.WriteString(" ")
		}

		hasAny = true

		fmt.Fprintf(
			sb,
			"%c%s%c%s%s",
			zittish.OpGroupOpen,
			ids,
			zittish.OpGroupClose,
			ids.Sigil,
			g,
		)
	}

	return sb.String()
}

func (s *group) BuildQueryGroup(vs ...string) (qg matcher_proto.QueryGroup, err error) {
	var builder query.Builder

	builder.
		WithDefaultGattungen(s.DefaultGattungen).
		WithCwd(s.cwd).
		WithFileExtensionGetter(s.fileExtensionGetter).
		WithHidden(s.Hidden).
		WithExpanders(s.expanders)

	if qg, err = builder.BuildQueryGroup(vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return

	if err = s.SetMany(vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg = s

	return
}

func (s *group) SetMany(vs ...string) (err error) {
	for _, v := range vs {
		if err = s.set(v); err != nil {
			var fd fd.FD

			if err1 := fd.Set(v); err1 == nil {
				if err = s.FDs.Add(&fd); err != nil {
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

func (ms *group) set(v string) (err error) {
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

	var gs kennung.Gattung

	if after != "" {
		if err = gs.Set(after); err != nil {
			if gattung.IsErrUnrecognizedGattung(err) {
				err = nil

				var f fd.FD

				if err = f.Set(v); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = ms.FDs.Add(&f); err != nil {
					err = errors.Wrap(err)
					return
				}

			} else {
				err = errors.Wrap(err)
			}

			return
		}
	} else {
		gs = ms.DefaultGattungen
	}

	if before == "" && after == "" && sigil.IncludesCwd() {
		ms.dotOperatorActive = true
		// return
	}

	for _, g := range gs.Slice() {
		var ids Query
		ok := false

		if ids, ok = ms.Gattung[g]; !ok {
			ids.Matcher = MakeMatcherExactlyThisOrAllOfThese()
			ids.Sigil = sigil
			ids.Gattung = g
		}

		switch {
		case before == "":
			break

		case ids.IncludesCwd():
			fp := fmt.Sprintf("%s.%s", before, after)

			var f *fd.FD

			if f, err = fd.FDFromPath(fp); err == nil {
				ids.Matcher.AddExactlyThis(FD{FD: f})
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

func (ms group) Get(g gattung.Gattung) (s MatcherSigil, ok bool) {
	var ids Query

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

func (ms group) GetExplicitCwdFDs() fd.Set {
	return ms.FDs
}

func (ms group) GetCwdFDs() fd.Set {
	if ms.dotOperatorActive {
		return ms.cwd.GetCwdFDs()
	} else {
		return ms.FDs
	}
}

func (ms group) GetEtiketten() kennung.EtikettSet {
	es := kennung.MakeEtikettMutableSet()

	for _, s := range ms.Gattung {
		VisitAllMatcherKennungSansGattungWrappers(
			func(m MatcherKennungSansGattungWrapper) (err error) {
				switch et := m.GetKennung().(type) {
				case kennung.Etikett:
					return es.AddPtr(&et)

				case *kennung.Etikett:
					return es.AddPtr(et)

				default:
					return
				}
			},
			IsMatcherNegate,
			// TODO-P1 modify sigil matcher to allow child traversal
			s.Matcher,
		)
	}

	return es
}

func (ms group) GetTyp() (t kennung.Typ, ok bool) {
	ts := ms.GetTypen()

	if ts.Len() != 1 {
		return
	}

	t = ts.Any()
	ok = true

	return
}

func (ms group) GetTypen() kennung.TypSet {
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

func (s group) ContainsMatchable(m *sku.Transacted) bool {
	g := gattung.Must(m.GetGattung())

	var matcher Matcher
	ok := false

	if matcher, ok = s.Get(g); !ok {
		return false
	}

	return matcher.ContainsMatchable(m)
}

func (ms group) GetGattungen() gattungen.Set {
	gs := make([]gattung.Gattung, 0, len(ms.Gattung))

	for g := range ms.Gattung {
		gs = append(gs, g)
	}

	return gattungen.MakeSet(gs...)
}
