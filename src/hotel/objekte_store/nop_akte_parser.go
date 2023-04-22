package objekte_store

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type nopAkteParser[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
] struct {
	arf schnittstellen.AkteIOFactory
}

func MakeNopAkteFormat[
	O objekte.Objekte[O],
	OPtr objekte.ObjektePtr[O],
](arf schnittstellen.AkteIOFactory,
) *nopAkteParser[O, OPtr] {
	return &nopAkteParser[O, OPtr]{
		arf: arf,
	}
}

func (f nopAkteParser[O, OPtr]) ParseAkte(
	r io.Reader,
	t OPtr,
) (n int64, err error) {
	var aw sha.WriteCloser

	if aw, err = f.arf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if n, err = io.Copy(aw, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.SetAkteSha(sha.Make(aw.Sha()))

	return
}

func (f nopAkteParser[O, OPtr]) Format(w io.Writer, t OPtr) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.GetAkteSha()); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.DeferredCloser(&err, ar)

	if n, err = io.Copy(w, ar); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
