package query

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/echo/zittish"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Exp struct {
	MatchOnEmpty bool
	Or           bool
	Negated      bool
	Exact        bool
	Hidden       bool
	Debug        bool
	Children     []sku.Query
}

func (a *Exp) Clone() (b *Exp) {
	b = &Exp{
		MatchOnEmpty: a.MatchOnEmpty,
		Or:           a.Or,
		Negated:      a.Negated,
		Exact:        a.Exact,
		Hidden:       a.Hidden,
		Debug:        a.Debug,
	}

	b.Children = make([]sku.Query, len(a.Children))

	for i, c := range a.Children {
		switch ct := c.(type) {
		case *Exp:
			b.Children[i] = ct.Clone()

		default:
			b.Children[i] = ct
		}
	}

	return b
}

func (e *Exp) CollectEtiketten(mes kennung.EtikettMutableSet) {
	if e.Or || e.Negated {
		return
	}

	for _, m := range e.Children {
		switch mt := m.(type) {
		case *Exp:
			mt.CollectEtiketten(mes)

		case *Kennung:
			if mt.Kennung2.GetGattung() != gattung.Etikett {
				continue
			}

			e := kennung.MustEtikett(mt.Kennung2.String())
			mes.Add(e)
		}
	}
}

func (e *Exp) Reduce(b *Builder) (err error) {
	e.MatchOnEmpty = !b.doNotMatchEmpty
	chillen := make([]sku.Query, 0, len(e.Children))

	for _, m := range e.Children {
		switch mt := m.(type) {
		case *Exp:
			if err = mt.Reduce(b); err != nil {
				err = errors.Wrap(err)
				return
			}

			if len(mt.Children) == 0 {
				continue
			}

			if mt.Or == e.Or && mt.Negated == e.Negated && mt.Exact == e.Exact {
				chillen = append(chillen, mt.Children...)
				continue
			}

		case Reducer:
			if err = mt.Reduce(b); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		chillen = append(chillen, m)
	}

	e.Children = chillen

	return
}

func (e *Exp) Add(m sku.Query) (err error) {
	switch mt := m.(type) {
	case *Exp:

	case *Kennung:
		mt.Exact = e.Exact
	}

	e.Children = append(e.Children, m)

	return
}

func (matcher *Exp) Operator() rune {
	if matcher.Or {
		return zittish.OpOr
	} else {
		return zittish.OpAnd
	}
}

func (e *Exp) StringDebug() string {
	var sb strings.Builder

	op := e.Operator()

	if e.Negated {
		sb.WriteRune('^')
	}

	if e.Exact {
		sb.WriteRune('=')
	}

	sb.WriteRune(zittish.OpGroupOpen)
	fmt.Fprintf(&sb, "(%d)", len(e.Children))

	for i, m := range e.Children {
		if i > 0 {
			sb.WriteRune(op)
		}

		sb.WriteString(m.String())
	}

	sb.WriteRune(zittish.OpGroupClose)

	return sb.String()
}

func (e *Exp) String() string {
	if e.Hidden {
		return ""
	}

	l := len(e.Children)

	if l == 0 {
		return ""
	}

	var sb strings.Builder

	op := e.Operator()

	if e.Negated {
		sb.WriteRune('^')
	}

	if e.Exact {
		sb.WriteRune('=')
	}

	switch l {
	case 1:
		sb.WriteString(e.Children[0].String())

	default:
		sb.WriteRune(zittish.OpGroupOpen)

		for i, m := range e.Children {
			if i > 0 {
				sb.WriteRune(op)
			}

			sb.WriteString(m.String())
		}

		sb.WriteRune(zittish.OpGroupClose)
	}

	return sb.String()
}

func (m *Exp) negateIfNecessary(v bool) bool {
	if m.Negated {
		return !v
	} else {
		return v
	}
}

func (e *Exp) ContainsSku(sk *sku.Transacted) bool {
	log.Log().Printf("%s in %s", sk, e)

	if len(e.Children) == 0 {
		return e.negateIfNecessary(e.MatchOnEmpty)
	}

	if e.Or {
		return e.containsMatchableOr(sk)
	} else {
		return e.containsMatchableAnd(sk)
	}
}

func (e *Exp) containsMatchableAnd(sk *sku.Transacted) bool {
	for _, m := range e.Children {
		if !m.ContainsSku(sk) {
			return e.negateIfNecessary(false)
		}
	}

	return e.negateIfNecessary(true)
}

func (e *Exp) containsMatchableOr(sk *sku.Transacted) bool {
	for _, m := range e.Children {
		if m.ContainsSku(sk) {
			return e.negateIfNecessary(true)
		}
	}

	return e.negateIfNecessary(false)
}

func (e *Exp) Each(
	f schnittstellen.FuncIter[sku.Query],
) (err error) {
	for _, m := range e.Children {
		if err = f(m); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
