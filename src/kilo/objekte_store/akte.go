package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type AkteStore[
	A schnittstellen.Akte[A],
	APtr schnittstellen.AktePtr[A],
] struct {
	standort standort.Standort
	StoredParseSaver[A, APtr]
	objekte.AkteFormat[A, APtr]
}

func MakeAkteStore[
	A schnittstellen.Akte[A],
	APtr schnittstellen.AktePtr[A],
](
	st standort.Standort,
	akteFormat objekte.AkteFormat[A, APtr],
) (s *AkteStore[A, APtr]) {
	s = &AkteStore[A, APtr]{
		standort:   st,
		AkteFormat: akteFormat,
	}

	return
}

func (s *AkteStore[A, APtr]) GetAkte(
	sh schnittstellen.ShaLike,
) (a APtr, err error) {
	var ar schnittstellen.ShaReadCloser

	if ar, err = s.standort.AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var a1 A
	a = APtr(&a1)
	a.Reset()

	if _, err = s.ParseAkte(ar, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	actual := ar.GetShaLike()

	if !actual.EqualsSha(sh) {
		err = errors.Errorf("expected sha %s but got %s", sh, actual)
		return
	}

	return
}

func (s *AkteStore[A, APtr]) PutAkte(a APtr) {
	// TODO-P2 implement pool
}

func (h *AkteStore[A, APtr]) SaveAkteText(
	o APtr,
) (sh schnittstellen.ShaLike, n int64, err error) {
	var w sha.WriteCloser

	if w, err = h.standort.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if n, err = h.AkteFormat.FormatParsedAkte(w, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.Make(w.GetShaLike())

	return
}
