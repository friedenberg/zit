package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type Element interface {
	schnittstellen.ValueLike
	schnittstellen.Stored
	AkteEtiketten() EtikettSet
	AkteTyp() Typ
	Hinweis() Hinweis
}

type Filter struct {
	Set Set
	Or  bool
}

// TODO-P4 improve the performance of this query
func (f Filter) Include(e Element) (err error) {
	ok := false

	// TODO-P3 pull into static
	needsEt := f.Set.Etiketten.Len() > 0
	okEt := false

	expanded := Expanded(e.AkteEtiketten(), ExpanderRight)

	// TODO-P3 pull into static
	ets := collections.SortedValues[Etikett](
		collections.Map[Etikett, Etikett](
			f.Set.Etiketten,
			func(e Etikett) (f Etikett) {
				return SansPrefix(e)
			},
		),
	)

LOOP:
	for _, e := range ets {
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

	// TODO-P2 make static
	shas := f.Set.Shas.ImmutableClone()
	needsSha := shas.Len() > 0
	okSha := false

	switch {
	case shas.Contains(sha.Make(e.GetObjekteSha())):
		okSha = true

		ok = false
		ok = false
		ok = false
	case shas.Contains(sha.Make(e.GetAkteSha())):
		okSha = true
	}

	ty := f.Set.Typen.ImmutableClone()
	needsTyp := ty.Len() > 0
	okTyp := false

	ty.Each(
		func(t Typ) (err error) {
			if okTyp = t.Includes(e.AkteTyp()); okTyp {
				err = collections.MakeErrStopIteration()
			}

			return
		},
	)

	hinweisen := f.Set.Hinweisen.ImmutableClone()
	needsHin := hinweisen.Len() > 0
	okHin := false || hinweisen.Len() == 0

	okHin = hinweisen.Contains(e.Hinweis())

	isEmpty := !needsHin && !needsTyp && !needsEt && !needsSha

	switch {
	case isEmpty:
		ok = f.Set.Sigil.IncludesSchwanzen() || f.Set.Sigil.IncludesHistory()

	case f.Or:
		ok = (okHin && needsHin) || (okTyp && needsTyp) || (okEt && needsEt) || (okSha && needsSha)

	default:
		ok = (okHin || !needsHin) && (okTyp || !needsTyp) && (okEt || !needsEt) && (okSha || !needsSha)
	}

	if !ok {
		err = collections.MakeErrStopIteration()
	}

	return
}
