package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/zittish"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
)

type Builder struct {
	implicitEtikettenGetter matcher_proto.ImplicitEtikettenGetter
	cwd                     matcher_proto.Cwd
	fileExtensionGetter     schnittstellen.FileExtensionGetter
	expanders               kennung.Abbr
	hidden                  matcher_proto.Matcher
	defaultGattungen        kennung.Gattung
	gattung                 map[kennung.Gattung]matcher_proto.MatcherExactlyThisOrAllOfThese
	doNotMatchEmpty         bool
	debug                   bool
}

func (mb *Builder) WithDebug() *Builder {
	mb.debug = true
	return mb
}

func (mb *Builder) WithCwd(
	cwd matcher_proto.Cwd,
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

func (mb *Builder) WithHidden(
	hidden matcher_proto.Matcher,
) *Builder {
	mb.hidden = hidden
	return mb
}

func (mb *Builder) WithImplicitEtikettenGetter(
	ieg matcher_proto.ImplicitEtikettenGetter,
) *Builder {
	mb.implicitEtikettenGetter = ieg
	return mb
}

func (b *Builder) BuildQueryGroup(vs ...string) (qg matcher_proto.QueryGroup, err error) {
	if qg, err = b.build(vs...); err != nil {
		err = errors.Wrapf(err, "Query: %q", vs)
		return
	}

	return
}

func (b *Builder) build(vs ...string) (qg *QueryGroup, err error) {
	qg = MakeQueryGroup(b)

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

	log.Log().Print(remaining)

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

			if err = qg.AddExactKennung(b, k); err != nil {
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

	log.Log().Print(qg)

	for _, q := range qg.OptimizedQueries {
		q.Reduce(b)
	}

	log.Log().Print(qg)

	return
}

func (b *Builder) buildManyFromTokens(
	qg *QueryGroup,
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

func (b *Builder) makeQuery() *Query {
	return &Query{
		Kennung: make(map[string]*kennung.Kennung2),
	}
}

func (b *Builder) makeExp(negated, exact bool, children ...matcher_proto.Matcher) *Exp {
	return &Exp{
		// MatchOnEmpty: !b.doNotMatchEmpty,
		Negated:  negated,
		Exact:    exact,
		Children: children,
	}
}

func (b *Builder) parseOneFromTokens(
	qg *QueryGroup,
	tokens ...string,
) (remainingTokens []string, err error) {
	q := b.makeQuery()
	stack := []matcher_proto.MatcherParentPtr{q}
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

			if gattung.Must(k) != gattung.Zettel {
				exp := b.makeExp(isNegated, isExact, &k)
				stack[len(stack)-1].Add(exp)
			}

			isNegated = false
			isExact = false

			switch k.GetGattung() {
			case gattung.Zettel:
				q.Kennung[k.String()] = k.Kennung2

			case gattung.Etikett:
				var e kennung.Etikett

				if err = e.TodoSetFromKennung2(k.Kennung2); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = qg.Etiketten.Add(e); err != nil {
					err = errors.Wrap(err)
					return
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
			}
		}
	}

	if q.IsEmpty() {
		q.Gattung = b.defaultGattungen
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

			q.Sigil.Add(s)
		}
	}

	if remainingTokens, err = q.SetTokens(remainingTokens...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
