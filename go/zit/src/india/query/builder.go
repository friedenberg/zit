package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/standort"
	"code.linenisgreat.com/zit/src/delta/zittish"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func MakeBuilder(s standort.Standort) *Builder {
	return &Builder{
		standort: s,
	}
}

type Builder struct {
	standort                standort.Standort
	preexistingKennung      []*kennung.Kennung2
	implicitEtikettenGetter ImplicitEtikettenGetter
	cwd                     Cwd
	fileExtensionGetter     schnittstellen.FileExtensionGetter
	expanders               kennung.Abbr
	hidden                  Matcher
	defaultGattungen        kennung.Gattung
	defaultSigil            kennung.Sigil
	doNotMatchEmpty         bool
	debug                   bool
}

func (mb *Builder) WithDebug() *Builder {
	mb.debug = true
	return mb
}

func (mb *Builder) WithCwd(
	cwd Cwd,
) *Builder {
	mb.cwd = cwd
	return mb
}

func (mb *Builder) WithFileExtensionGetter(
	feg schnittstellen.FileExtensionGetter,
) *Builder {
	mb.fileExtensionGetter = feg
	return mb
}

func (mb *Builder) WithExpanders(
	expanders kennung.Abbr,
) *Builder {
	mb.expanders = expanders
	return mb
}

func (mb *Builder) WithDefaultGattungen(
	defaultGattungen kennung.Gattung,
) *Builder {
	mb.defaultGattungen = defaultGattungen
	return mb
}

func (mb *Builder) WithDefaultSigil(
	defaultSigil kennung.Sigil,
) *Builder {
	mb.defaultSigil = defaultSigil
	return mb
}

func (mb *Builder) WithHidden(
	hidden Matcher,
) *Builder {
	mb.hidden = hidden
	return mb
}

func (mb *Builder) WithImplicitEtikettenGetter(
	ieg ImplicitEtikettenGetter,
) *Builder {
	mb.implicitEtikettenGetter = ieg
	return mb
}

func (b *Builder) WithTransacted(
	zts sku.TransactedSet,
) *Builder {
	errors.PanicIfError(zts.Each(
		func(t *sku.Transacted) (err error) {
			k := kennung.GetKennungPool().Get()
			k.ResetWith(&t.Kennung)
			b.preexistingKennung = append(b.preexistingKennung, k)

			return
		},
	))

	return b
}

func (b *Builder) WithCheckedOut(
	cos sku.CheckedOutSet,
) *Builder {
	errors.PanicIfError(cos.Each(
		func(co *sku.CheckedOut) (err error) {
			k := kennung.GetKennungPool().Get()
			k.ResetWith(&co.Internal.Kennung)
			b.preexistingKennung = append(b.preexistingKennung, k)

			return
		},
	))

	return b
}

func (b *Builder) BuildQueryGroup(vs ...string) (qg *Group, err error) {
	if qg, err = b.build(vs...); err != nil {
		err = errors.Wrapf(err, "Query: %q", vs)
		return
	}

	log.Log().Print(qg.StringDebug())

	return
}

func (b *Builder) build(vs ...string) (qg *Group, err error) {
	qg = MakeGroup(b)

	var remaining []string

	for _, v := range vs {
		var fd fd.FD

		if err1 := fd.Set(v); err1 != nil {
			remaining = append(remaining, v)
			continue
		}

		if err = qg.FDs.Add(&fd); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var tokens []string

	if tokens, err = zittish.GetTokensFromStrings(remaining...); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = b.buildManyFromTokens(qg, tokens...); err != nil {
		err = errors.Wrap(err)
		return
	}

	newFDS := fd.MakeMutableSet()

	if err = qg.FDs.Each(
		func(f *fd.FD) (err error) {
			var k *kennung.Kennung2

			if k, err = b.cwd.GetKennungForFD(f); err != nil {
				if errors.Is(err, kennung.ErrFDNotKennung) {
					if err = newFDS.Add(f); err != nil {
						err = errors.Wrap(err)
						return
					}
				} else {
					err = errors.Wrapf(err, "File: %q", f)
				}

				return
			}

			if err = qg.AddExactKennung(
				b,
				Kennung{Kennung2: k, FD: f},
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	qg.FDs = newFDS

	for _, k := range b.preexistingKennung {
		if err = qg.AddExactKennung(
			b,
			Kennung{Kennung2: k},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	b.addDefaultsIfNecessary(qg)

	if err = qg.Reduce(b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *Builder) buildManyFromTokens(
	qg *Group,
	tokens ...string,
) (err error) {
	if len(tokens) == 1 && tokens[0] == "." {
		if err = b.cwd.GetCwdFDs().Each(qg.FDs.Add); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	for len(tokens) > 0 {
		if tokens, err = b.parseOneFromTokens(qg, tokens...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (b *Builder) addDefaultsIfNecessary(qg *Group) {
	if b.defaultGattungen.IsEmpty() || !qg.IsEmpty() {
		return
	}

	g := kennung.MakeGattung()
	dq, ok := qg.UserQueries[g]

	if ok {
		delete(qg.UserQueries, g)
	} else {
		dq = &QueryWithHidden{
			Query: Query{
				Kennung: make(map[string]Kennung),
			},
		}
	}

	dq.Gattung = b.defaultGattungen

	if b.defaultSigil.IsEmpty() {
		dq.Sigil = kennung.SigilSchwanzen
	} else {
		dq.Sigil = b.defaultSigil
	}

	qg.UserQueries[b.defaultGattungen] = dq

	return
}

func (b *Builder) makeQuery() *Query {
	return &Query{
		Kennung: make(map[string]Kennung),
	}
}

func (b *Builder) makeExp(negated, exact bool, children ...Matcher) *Exp {
	return &Exp{
		// MatchOnEmpty: !b.doNotMatchEmpty,
		Negated:  negated,
		Exact:    exact,
		Children: children,
	}
}

func (b *Builder) parseOneFromTokens(
	qg *Group,
	tokens ...string,
) (remainingTokens []string, err error) {
	type stackEl interface {
		Matcher
		Add(Matcher) error
	}

	q := b.makeQuery()
	stack := []stackEl{q}

	isNegated := false
	isExact := false

LOOP:
	for i, el := range tokens {
		if len(el) == 1 && zittish.IsMatcherOperator([]rune(el)[0]) {
			op := el[0]
			switch op {
			case '=':
				isExact = true

			case '^':
				isNegated = true

			case ' ':

			case ',':
				last := stack[len(stack)-1].(*Exp)
				last.Or = true
				// TODO handle or when invalid

			case '[':
				exp := b.makeExp(isNegated, isExact)
				isExact = false
				isNegated = false
				stack[len(stack)-1].Add(exp)
				stack = append(stack, exp)

			case ']':
				stack = stack[:len(stack)-1]
				// TODO handle errors of unbalanced

			case '.':
				// TODO end sigil or embedded as part of name
				fallthrough

			case ':', '+', '?':
				if len(stack) > 1 {
					err = errors.Errorf("sigil before end")
					return
				}

				if remainingTokens, err = b.parseSigilsAndGattungen(q, tokens[i:]...); err != nil {
					err = errors.Wrapf(err, "%s", tokens[i:])
					return
				}

				break LOOP
			}
		} else {
			k := Kennung{
				Kennung2: kennung.GetKennungPool().Get(),
			}

			if err = k.Set(el); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = k.Reduce(b); err != nil {
				err = errors.Wrap(err)
				return
			}

			switch k.GetGattung() {
			case gattung.Zettel:
				b.preexistingKennung = append(b.preexistingKennung, k.Kennung2)
				q.Gattung.Add(gattung.Zettel)
				q.Kennung[k.Kennung2.String()] = k

			case gattung.Etikett:
				var e kennung.Etikett

				if err = e.TodoSetFromKennung2(k.Kennung2); err != nil {
					err = errors.Wrap(err)
					return
				}

				if e.IsVirtual() && strings.HasPrefix(e.String(), "%chrome") {
					c := MakeChrome(b.standort)

					if err = c.Init(); err != nil {
						err = errors.Wrap(err)
						return
					}

					exp := b.makeExp(isNegated, isExact, c)
					stack[len(stack)-1].Add(exp)
				} else {
					if err = qg.Etiketten.Add(e); err != nil {
						err = errors.Wrap(err)
						return
					}

					exp := b.makeExp(isNegated, isExact, &k)
					stack[len(stack)-1].Add(exp)
				}

			case gattung.Typ:
				var t kennung.Typ

				if err = t.TodoSetFromKennung2(k.Kennung2); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = qg.Typen.Add(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				exp := b.makeExp(isNegated, isExact, &k)
				stack[len(stack)-1].Add(exp)
			}

			isNegated = false
			isExact = false
		}
	}

	if q.IsEmpty() {
		return
	}

	if q.Gattung.IsEmpty() {
		q.Gattung = b.defaultGattungen
	}

	if q.Sigil.IsEmpty() {
		q.Sigil = b.defaultSigil
	}

	if err = qg.Add(q); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *Builder) parseSigilsAndGattungen(
	q *Query,
	tokens ...string,
) (remainingTokens []string, err error) {
LOOP:
	for i, el := range tokens {
		if len(el) != 1 {
			remainingTokens = tokens[i:]
			break
		}

		op := []rune(el)[0]

		switch op {
		default:
			remainingTokens = tokens[i:]
			break LOOP

		case ':', '+', '?', '.':
			var s kennung.Sigil

			if err = s.Set(el); err != nil {
				err = errors.Wrap(err)
				return
			}

			if op == '?' {
				q.Sigil.Add(kennung.SigilSchwanzen)
			}

			q.Sigil.Add(s)
		}
	}

	if remainingTokens, err = q.SetTokens(remainingTokens...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
