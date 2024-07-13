package akten

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
)

type AkteStore[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] struct {
	standort standort.Standort
	Format[A, APtr]
	resetFunc func(APtr)
}

func MakeAkteStore[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
](
	st standort.Standort,
	akteFormat Format[A, APtr],
	resetFunc func(APtr),
) (s *AkteStore[A, APtr]) {
	s = &AkteStore[A, APtr]{
		standort:  st,
		Format:    akteFormat,
		resetFunc: resetFunc,
	}

	return
}

func (s *AkteStore[A, APtr]) GetBlob(
	sh interfaces.ShaLike,
) (a APtr, err error) {
	var ar interfaces.ShaReadCloser

	if ar, err = s.standort.BlobReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ar)

	var a1 A
	a = APtr(&a1)
	s.resetFunc(a)

	if _, err = s.ParseBlob(ar, a); err != nil {
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

func (s *AkteStore[A, APtr]) PutBlob(a APtr) {
	// TODO-P2 implement pool
}

func (h *AkteStore[A, APtr]) SaveAkteText(
	o APtr,
) (sh interfaces.ShaLike, n int64, err error) {
	var w sha.WriteCloser

	if w, err = h.standort.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, w)

	if n, err = h.FormatParsedAkte(w, o); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = w.GetShaLike()

	return
}
