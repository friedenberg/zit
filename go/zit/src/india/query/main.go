package query

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type Query struct {
	kennung.Sigil
	kennung.Gattung
	Exp

	Kennung map[string]*kennung.Kennung2
}

func (a *Query) GetMatcherSigil() matcher_proto.MatcherSigil {
	return a
}

func (a *Query) GetSigil() kennung.Sigil {
	return a.Sigil
}

func (a *Query) Clone() (b *Query) {
	b = &Query{
		Sigil:   a.Sigil,
		Gattung: a.Gattung,
		Kennung: make(map[string]*kennung.Kennung2, len(a.Kennung)),
	}

	bExp := a.Exp.Clone()
	b.Exp = *bExp

	for k, v := range a.Kennung {
		b.Kennung[k] = v
	}

	return b
}

func (q *Query) Add(m matcher_proto.Matcher) (err error) {
	log.Log().Print("before", q.StringDebug(), m)
	defer log.Log().Print("after", q.StringDebug())
	q1, ok := m.(*Query)

	if !ok {
		return q.Exp.Add(m)
	}

	if q1.Gattung != q.Gattung {
		err = errors.Errorf(
			"expected %q but got %q",
			q.Gattung,
			q1.Gattung,
		)

		return
	}

	if err = q.Merge(q1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Query) Merge(b *Query) (err error) {
	a.Sigil.Add(b.Sigil)

	if a.Kennung == nil {
		a.Kennung = make(map[string]*kennung.Kennung2)
	}

	for _, k := range b.Kennung {
		a.Kennung[k.String()] = k
	}

	if err = a.Exp.Add(&b.Exp); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (q *Query) Each(_ schnittstellen.FuncIter[matcher_proto.Matcher]) (err error) {
	return
}

func (q *Query) MatcherLen() int {
	return 0
}

func (q *Query) StringDebug() string {
	var sb strings.Builder

	if q.Kennung == nil || len(q.Kennung) == 0 {
		sb.WriteString(q.Exp.StringDebug())
	} else {
		sb.WriteString("[[")

		first := true

		for _, k := range q.Kennung {
			if !first {
				sb.WriteString(", ")
			}

			sb.WriteString(k.String())

			first = false
		}

		sb.WriteString(", ")
		sb.WriteString(q.Exp.StringDebug())
		sb.WriteString("]")
	}

	if q.IsEmpty() && !q.IsSchwanzenOrUnknown() {
		sb.WriteString(q.Sigil.String())
	} else if !q.IsEmpty() {
		sb.WriteString(q.Sigil.String())
		sb.WriteString(q.Gattung.String())
	}

	return sb.String()
}

func (q *Query) String() string {
	var sb strings.Builder

	e := q.Exp.String()
	log.Log().Printf("%q", e)

	if q.Kennung == nil || len(q.Kennung) == 0 {
		sb.WriteString(e)
	} else if len(q.Kennung) == 1 && e == "" {
		for _, k := range q.Kennung {
			sb.WriteString(k.String())
		}
	} else {
		sb.WriteString("[")

		first := true

		for _, k := range q.Kennung {
			if !first {
				sb.WriteString(", ")
			}

			sb.WriteString(k.String())

			first = false
		}

		if e != "" {
			sb.WriteString(", ")
			sb.WriteString(q.Exp.String())
		}

		sb.WriteString("]")
	}

	if q.IsEmpty() && !q.IsSchwanzenOrUnknown() {
		sb.WriteString(q.Sigil.String())
	} else if !q.IsEmpty() {
		sb.WriteString(q.Sigil.String())
		sb.WriteString(q.Gattung.String())
	}

	return sb.String()
}

func (q *Query) ContainsMatchable(sk *sku.Transacted) bool {
	log.Log().Print(q, sk)
	g := gattung.Must(sk)

	if !q.Gattung.Contains(g) {
		log.Log().Print(q, sk, g)
		return false
	}

	log.Log().Print(q, q.Kennung, sk.Kennung.String())

	if _, ok := q.Kennung[sk.Kennung.String()]; ok {
		log.Log().Print(q, sk)
		return true
	}

	log.Log().Printf("%#v", q.Children)
	log.Log().Print(q, sk, len(q.Children), q.Children, len(q.Kennung))

	if len(q.Children) == 0 {
		log.Log().Print(q, sk, len(q.Kennung))
		return len(q.Kennung) == 0
	} else {
		if !q.Exp.ContainsMatchable(sk) {
			return false
		}
	}

	return true
}
