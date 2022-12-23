package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
)

type hydrator[T gattung.Element, T1 gattung.ElementPtr[T]] struct {
	arf              gattung.AkteReaderFactory
	frc              FuncReadCloser
	objekteFormatter Formatter2
	akteParser       gattung.Parser[T, T1] //TODO-P1 rename to akteParser
}

func MakeHydrator[T gattung.Element, T1 gattung.ElementPtr[T]](
	arf gattung.AkteReaderFactory,
	frc FuncReadCloser,
	akteParser gattung.Parser[T, T1],
) *hydrator[T, T1] {
	return &hydrator[T, T1]{
		arf:              arf,
		frc:              frc,
		objekteFormatter: *MakeFormatter2(),
		akteParser:       akteParser,
	}
}

func (h *hydrator[T, T1]) Hydrate(
	to gattung.StoredPtr,
	a T1,
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

	if h.akteParser != nil {
		var r sha.ReadCloser

		if r, err = h.arf.AkteReader(to.AkteSha()); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.Deferred(&err, r.Close)

		if _, err = h.akteParser.Parse(r, a); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
