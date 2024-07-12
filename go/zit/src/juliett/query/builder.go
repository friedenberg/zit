package query

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/etikett_akte"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/akten"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func MakeBuilder(
	s standort.Standort,
	akten *akten.Akten,
	ennui sku.Ennui,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	kastenGetter sku.ExternalStoreGetter,
) (b *Builder) {
	b = &Builder{
		standort:                   s,
		akten:                      akten,
		ennui:                      ennui,
		luaVMPoolBuilder:           luaVMPoolBuilder,
		virtualEtikettenBeforeInit: make(map[string]string),
		virtualEtiketten:           make(map[string]Lua),
		kastenGetter:               kastenGetter,
	}

	return
}

type Builder struct {
	standort                   standort.Standort
	akten                      *akten.Akten
	ennui                      sku.Ennui
	luaVMPoolBuilder           *lua.VMPoolBuilder
	preexistingKennung         []Kennung
	kastenGetter               sku.ExternalStoreGetter
	kasten                     *sku.ExternalStore
	cwdFilterEnabled           bool
	fileExtensionGetter        schnittstellen.FileExtensionGetter
	expanders                  kennung.Abbr
	hidden                     sku.Query
	defaultGattungen           kennung.Gattung
	defaultSigil               kennung.Sigil
	permittedSigil             kennung.Sigil
	virtualEtikettenBeforeInit map[string]string
	virtualEtiketten           map[string]Lua
	doNotMatchEmpty            bool
	debug                      bool
	requireNonEmptyQuery       bool
	eqo                        sku.ExternalQueryOptions
}

func (b *Builder) WithPermittedSigil(s kennung.Sigil) *Builder {
	b.permittedSigil.Add(s)
	return b
}

func (b *Builder) WithDoNotMatchEmpty() *Builder {
	b.doNotMatchEmpty = true
	return b
}

func (b *Builder) WithCwdFilterEnabled() *Builder {
	b.cwdFilterEnabled = true
	return b
}

func (b *Builder) WithRequireNonEmptyQuery() *Builder {
	b.requireNonEmptyQuery = true
	return b
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

func (mb *Builder) WithKasten(
	kasten *sku.ExternalStore,
) *Builder {
	mb.kasten = kasten
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

func (b *Builder) WithTransacted(
	zts sku.TransactedSet,
) *Builder {
	errors.PanicIfError(zts.Each(
		func(t *sku.Transacted) (err error) {
			b.preexistingKennung = append(
				b.preexistingKennung,
				Kennung{
					Kennung2: t.Kennung.Clone(),
				},
			)

			return
		},
	))

	return b
}

func (b *Builder) WithCheckedOut(
	cos sku.CheckedOutLikeSet,
) *Builder {
	errors.PanicIfError(cos.Each(
		func(co sku.CheckedOutLike) (err error) {
			b.preexistingKennung = append(
				b.preexistingKennung,
				Kennung{
					Kennung2: co.GetSku().Kennung.Clone(),
				},
			)

			return
		},
	))

	return b
}

func (b *Builder) realizeVirtualEtiketten() (err error) {
	for k, v := range b.virtualEtikettenBeforeInit {
		var vmp *lua.VMPool

		lb := b.luaVMPoolBuilder.Clone().WithScript(v)

		if vmp, err = lb.Build(); err != nil {
			err = errors.Wrap(err)
			return
		}

		ml := Lua{
			LuaVMPool: sku_fmt.MakeLuaVMPool(vmp, nil),
		}

		b.virtualEtiketten[k] = ml
	}

	return
}

func (b *Builder) BuildQueryGroupWithKasten(
	k kennung.Kasten,
	eqo sku.ExternalQueryOptions,
	vs ...string,
) (qg *Group, err error) {
	ok := false
	b.eqo = eqo
	b.kasten, ok = b.kastenGetter.GetExternalStore(k)

	if !ok {
		err = errors.Errorf("kasten not found: %q", k)
		return
	}

	if qg, err = b.BuildQueryGroup(vs...); err != nil {
		err = errors.Wrap(err)
		return
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

	ui.Log().Print(qg.StringDebug())

	return
}

func (b *Builder) build(vs ...string) (qg *Group, err error) {
	qg = MakeGroup(b)

	var remaining []string

	for _, v := range vs {
		var k *kennung.Kennung2

		if k, err = b.kasten.GetKennungForString(v); err != nil {
			err = nil
			remaining = append(remaining, v)
			continue
		}

		b.preexistingKennung = append(
			b.preexistingKennung,
			Kennung{
				Kennung2: k,
				External: true,
			},
		)
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

	for _, k := range b.preexistingKennung {
		if err = qg.AddExactKennung(
			b,
			k,
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
		// TODO [ces/mew] switch to marker on query group for Cwd
		var ks schnittstellen.SetLike[*kennung.Kennung2]

		if ks, err = b.kasten.GetExternalKennung(); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = ks.Each(
			func(k *kennung.Kennung2) (err error) {
				b.preexistingKennung = append(
					b.preexistingKennung,
					Kennung{
						Kennung2: k,
						External: true,
					},
				)

				return
			},
		); err != nil {
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
		dq = &Query{
			Kennung: make(map[string]Kennung),
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
				b.preexistingKennung = append(
					b.preexistingKennung,
					k,
				)

				q.Gattung.Add(gattung.Zettel)
				q.Kennung[k.Kennung2.String()] = k

			case gattung.Etikett:
				var et sku.Query

				if et, err = b.makeEtikettExp(&k); err != nil {
					err = errors.Wrap(err)
					return
				}

				exp := b.makeExp(isNegated, isExact, et)
				stack[len(stack)-1].Add(exp)

			case gattung.Typ:
				var t kennung.Typ

				if err = t.TodoSetFromKennung2(k.Kennung2); err != nil {
					err = errors.Wrap(err)
					return
				}

				if !isNegated {
					if err = qg.Typen.Add(t); err != nil {
						err = errors.Wrap(err)
						return
					}
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

func (b *Builder) makeEtikettOrEtikettLua(
	k *Kennung,
) (exp sku.Query, err error) {
	exp = k

	if b.ennui == nil || b.akten == nil {
		return
	}

	var sk *sku.Transacted

	if sk, err = b.ennui.ReadOneKennung(k); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer sku.GetTransactedPool().Put(sk)

	lb := b.luaVMPoolBuilder.Clone().WithApply(MakeSelbstApply(sk))

	// TODO use repo pattern
	if sk.GetTyp().String() == "lua" {
		var ar sha.ReadCloser

		if ar, err = b.standort.AkteReader(sk.GetAkteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, ar)

		lb.WithReader(ar)
	} else {
		var akte *etikett_akte.V1

		if akte, err = b.akten.GetEtikettV1().GetAkte(
			sk.GetAkteSha(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if akte.Filter == "" {
			return
		}

		lb.WithScript(akte.Filter)
	}

	var vmp *lua.VMPool

	if vmp, err = lb.Build(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ml := Lua{
		LuaVMPool: sku_fmt.MakeLuaVMPool(vmp, nil),
	}

	exp = &EtikettLua{Lua: &ml, Kennung: k}

	return
}

func (b *Builder) makeEtikettExp(k *Kennung) (exp sku.Query, err error) {
	// TODO use b.akten to read Etikett Akte and find filter if necessary
	var e kennung.Etikett

	if err = e.TodoSetFromKennung2(k.Kennung2); err != nil {
		err = errors.Wrap(err)
		return
	}

	if exp, err = b.makeEtikettOrEtikettLua(k); err != nil {
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

			// if op == '?' {
			// 	q.Sigil.Add(kennung.SigilSchwanzen)
			// }

			q.Sigil.Add(s)
		}
	}

	if remainingTokens, err = q.SetTokens(remainingTokens...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
