package zettel_named

import (
	"io"

	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/foxtrot/id_set"
)

type FilterIdSet struct {
	id_set.Set
	AllowEmpty bool
	Or         bool
}

func (f FilterIdSet) WriteZettelNamed(z *Zettel) (err error) {
	if !f.IncludeNamedZettel(z) {
		err = io.EOF
		return
	}

	return
}

// TODO improve the performance of this query
func (f FilterIdSet) IncludeNamedZettel(z *Zettel) (ok bool) {
	needsEt := f.Set.Etiketten().Len() > 0
	okEt := false

	expanded := etikett.Expanded(z.Stored.Zettel.Etiketten, etikett.ExpanderRight{})

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
	case shas.Contains(z.Stored.Sha):
		okSha = true

	case shas.Contains(z.Stored.Zettel.Akte):
		okSha = true
	}

	needsTyp := len(f.Set.Typen()) > 0
	okTyp := false

	for _, t := range f.Set.Typen() {
		if okTyp = t.Includes(z.Stored.Zettel.Typ.Etikett); okTyp {
			break
		}
	}

	hinweisen := f.Set.Hinweisen()
	needsHin := hinweisen.Len() > 0
	okHin := false || hinweisen.Len() == 0

	okHin = hinweisen.Contains(z.Hinweis)

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

	// errors.Err().Print(okHin, needsHin, okTyp, needsTyp, okEt, needsEt, okSha, needsSha)

	return
}