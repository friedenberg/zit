package query

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/standort"
	"code.linenisgreat.com/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func MakeBuilder(
	s standort.Standort,
	chrome *VirtualStoreInitable,
) (b *Builder) {
	b = &Builder{
		standort:                   s,
		virtualStores:              make(map[string]*VirtualStoreInitable),
		virtualEtikettenBeforeInit: make(map[string]string),
		virtualEtiketten:           make(map[string]*Lua),
	}

	if chrome != nil {
		b.WithChrome(chrome)
	}

	return
}

type Builder struct {
	standort                   standort.Standort
	preexistingKennung         []*kennung.Kennung2
	implicitEtikettenGetter    ImplicitEtikettenGetter
	cwd                        Cwd
	fileExtensionGetter        schnittstellen.FileExtensionGetter
	expanders                  kennung.Abbr
	hidden                     sku.Query
	defaultGattungen           kennung.Gattung
	defaultSigil               kennung.Sigil
	permittedSigil             kennung.Sigil
	virtualStores              map[string]*VirtualStoreInitable
	virtualEtikettenBeforeInit map[string]string
	virtualEtiketten           map[string]*Lua
	doNotMatchEmpty            bool
	debug                      bool
	requireNonEmptyQuery       bool
}

func (b *Builder) WithPermittedSigil(s kennung.Sigil) *Builder {
	b.permittedSigil.Add(s)
	return b
}

func (b *Builder) WithDoNotMatchEmpty() *Builder {
	b.doNotMatchEmpty = true
	return b
}

func (b *Builder) WithRequireNonEmptyQuery() *Builder {
	b.requireNonEmptyQuery = true
	return b
}

func (mb *Builder) WithChrome(vs *VirtualStoreInitable) *Builder {
	mb.virtualStores["%chrome"] = vs

	return mb
}

func (mb *Builder) WithVirtualStores(vs map[string]*VirtualStoreInitable) *Builder {
	for k, v := range vs {
		mb.virtualStores[k] = v
	}

	return mb
}

func (mb *Builder) WithVirtualEtiketten(vs map[string]string) *Builder {
	for k, v := range vs {
		mb.virtualEtikettenBeforeInit["%"+k] = v
	}

	return mb
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
	hidden sku.Query,
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

func (b *Builder) realizeVirtualEtiketten() (err error) {
	for k, v := range b.virtualEtikettenBeforeInit {
		var ml *Lua

		if ml, err = MakeLua(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		b.virtualEtiketten[k] = ml
	}

	return
}

func (b *Builder) BuildQueryGroup(vs ...string) (qg *Group, err error) {
	if err = b.realizeVirtualEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	if b.requireNonEmptyQuery && qg.IsEmpty() {
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
}

func (b *Builder) makeQuery() *Query {
	return &Query{
		Kennung: make(map[string]Kennung),
	}
}

func (b *Builder) makeExp(negated, exact bool, children ...sku.Query) *Exp {
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
		sku.Query
		Add(sku.Query) error
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

				if e.IsVirtual() {
					expanded := kennung.ExpandOneSlice(&e)
					var store *VirtualStoreInitable
					var eStore *Lua

					for _, e1 := range expanded {
						store = b.virtualStores[e1.String()]

						if store != nil {
							break
						}

						eStore = b.virtualEtiketten[e1.String()]

						if eStore != nil {
							break
						}
					}

					if store == nil && eStore == nil {
						err = errors.Errorf("no virtual store registered for %q", e)
						return
					}

					if store != nil {
						if err = store.Initialize(); err != nil {
							err = errors.Wrap(err)
							return
						}

						exp := b.makeExp(isNegated, isExact, &Virtual{Queryable: store, Kennung: k})
						stack[len(stack)-1].Add(exp)
					} else {
						exp := b.makeExp(isNegated, isExact, &Virtual{Queryable: eStore, Kennung: k})
						stack[len(stack)-1].Add(exp)
					}
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

	if q.Gattung.IsEmpty() && !b.requireNonEmptyQuery {
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

			if !b.permittedSigil.IsEmpty() && !b.permittedSigil.ContainsOneOf(s) {
				err = errors.Errorf("cannot contain sigil %s", s)
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
