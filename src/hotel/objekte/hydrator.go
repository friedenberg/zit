package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

// type FuncReadCloser func(sha.Sha) (sha.ReadCloser, error)
// type FuncWriteCloser func(sha.Sha) (sha.WriteCloser, error)

type hydrator[T gattung.Element, T1 gattung.ElementPtr[T]] struct {
	af               gattung.AkteIOFactory
	frc              FuncReadCloser
	objekteFormatter Formatter2
	akteFormatter    gattung.Parser[T, T1] //TODO-P1 rename to akteParser
}

func MakeHydrator[T gattung.Element, T1 gattung.ElementPtr[T]](
	af gattung.AkteIOFactory,
	frc FuncReadCloser,
	akteFormatter gattung.Parser[T, T1],
) *hydrator[T, T1] {
	return &hydrator[T, T1]{
		af:               af,
		frc:              frc,
		objekteFormatter: *MakeFormatter2(),
		akteFormatter:    akteFormatter,
	}
}

func (h *hydrator[T, T1]) Hydrate(
	to gattung.StoredPtr,
	a *T,
) (err error) {
	{
		var r sha.ReadCloser

		if r, err = h.frc(to.ObjekteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, r.Close)

		if _, err = h.objekteFormatter.ReadFormat(r, to); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if h.akteFormatter != nil {
		var r sha.ReadCloser

		if r, err = h.af.AkteReader(to.AkteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, r.Close)

		if _, err = h.akteFormatter.Parse(r, a); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
