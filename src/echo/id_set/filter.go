package id_set

import (
	"io"

	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
)

type Element interface {
	gattung.Stored
	AkteEtiketten() kennung.EtikettSet
	AkteTyp() kennung.Typ
	Hinweis() hinweis.Hinweis
}

type Filter struct {
	Set
	AllowEmpty bool
	Or         bool
}

// TODO-P4 improve the performance of this query
func (f Filter) Include(e Element) (err error) {
	ok := false

	needsEt := f.Set.Etiketten().Len() > 0
	okEt := false

	expanded := kennung.Expanded(e.AkteEtiketten(), kennung.ExpanderRight)

LOOP:
	for _, e := range f.Set.Etiketten().Sorted() {
		okEt = expanded.Contains(e)

		switch {
		case !okEt && !f.Or:
			break LOOP

		case okEt && f.Or:
			break LOOP

		default:
			continue
		}
	}

	shas := f.Set.Shas()
	needsSha := shas.Len() > 0
	okSha := false

	switch {
	case shas.Contains(e.ObjekteSha()):
		okSha = true

	case shas.Contains(e.AkteSha()):
		okSha = true
	}

	needsTyp := len(f.Set.Typen()) > 0
	okTyp := false

	for _, t := range f.Set.Typen() {
		if okTyp = t.Includes(e.AkteTyp()); okTyp {
			break
		}
	}

	hinweisen := f.Set.Hinweisen()
	needsHin := hinweisen.Len() > 0
	okHin := false || hinweisen.Len() == 0

	okHin = hinweisen.Contains(e.Hinweis())

	isEmpty := !needsHin && !needsTyp && !needsEt && !needsSha

	switch {
	case isEmpty && f.AllowEmpty:
		ok = true

	case isEmpty:
		ok = false

	case f.Or:
		ok = (okHin && needsHin) || (okTyp && needsTyp) || (okEt && needsEt) || (okSha && needsSha)

	default:
		ok = (okHin || !needsHin) && (okTyp || !needsTyp) && (okEt || !needsEt) && (okSha || !needsSha)
	}

	if !ok {
		err = io.EOF
	}

	return
}
