package id_set

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/hinweis"
)

type Element interface {
	schnittstellen.Stored
	AkteEtiketten() kennung.EtikettSet
	AkteTyp() kennung.Typ
	Hinweis() hinweis.Hinweis
}

type Filter struct {
	Set        Set
	AllowEmpty bool
	Or         bool
}

// TODO-P4 improve the performance of this query
func (f Filter) Include(e Element) (err error) {
	ok := false

	//TODO-P3 pull into static
	needsEt := f.Set.Etiketten.Len() > 0
	okEt := false

	expanded := kennung.Expanded(e.AkteEtiketten(), kennung.ExpanderRight)

LOOP:
	//TODO-P3 pull into static
	for _, e := range f.Set.Etiketten.Copy().Sorted() {
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

	//TODO-P2 make static
	shas := f.Set.Shas.Copy()
	needsSha := shas.Len() > 0
	okSha := false

	switch {
	case shas.Contains(sha.Make(e.GetObjekteSha())):
		okSha = true

	case shas.Contains(sha.Make(e.GetAkteSha())):
		okSha = true
	}

	ty := f.Set.Typen.Copy()
	needsTyp := ty.Len() > 0
	okTyp := false

	ty.Each(
		func(t kennung.Typ) (err error) {
			if okTyp = t.Includes(e.AkteTyp()); okTyp {
				err = collections.ErrStopIteration
			}

			return
		},
	)

	hinweisen := f.Set.Hinweisen.Copy()
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
		err = collections.ErrStopIteration
	}

	return
}