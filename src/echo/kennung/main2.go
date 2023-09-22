package kennung

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type Kennung2 struct {
	KennungPtr
}

func MakeKennung2(v string) (KennungPtr, error) {
	k := &Kennung2{}
	return k, k.Set(v)
}

func (h *Kennung2) SetWithKennung(
	k Kennung,
) (err error) {
	h.KennungPtr, err = MakeWithGattung(k.GetGattung(), k.String())
	return
}

func (h *Kennung2) SetWithGattung(
	v string,
	g schnittstellen.GattungLike,
) (err error) {
	h.KennungPtr, err = MakeWithGattung(g, v)
	return
}

func (h *Kennung2) Set(v string) (err error) {
	h.KennungPtr, err = Make(v)
	return
}

func (a *Kennung2) ResetWithKennung(b Kennung) (err error) {
	switch bt := b.(type) {
	case Hinweis:
		a.KennungPtr = &bt

	case *Hinweis:
		b1 := *bt
		a.KennungPtr = &b1

	case Typ:
		a.KennungPtr = &bt

	case *Typ:
		b1 := *bt
		a.KennungPtr = &b1

	case Etikett:
		a.KennungPtr = &bt

	case *Etikett:
		b1 := *bt
		a.KennungPtr = &b1

	case Kasten:
		a.KennungPtr = &bt

	case *Kasten:
		b1 := *bt
		a.KennungPtr = &b1

	case Konfig:
		a.KennungPtr = &bt

	case *Konfig:
		b1 := *bt
		a.KennungPtr = &b1

	case Tai:
		a.KennungPtr = &bt

	case *Tai:
		b1 := *bt
		a.KennungPtr = &b1

	default:
		err = errors.Errorf("unsupported kennung: %T", b)
		return
	}

	return
}

func (a *Kennung2) ResetWithKennungPtr(b KennungPtr) (err error) {
	switch bt := b.(type) {
	case *Hinweis:
		b1 := *bt
		a.KennungPtr = &b1

	case *Typ:
		b1 := *bt
		a.KennungPtr = &b1

	case *Etikett:
		b1 := *bt
		a.KennungPtr = &b1

	case *Kasten:
		b1 := *bt
		a.KennungPtr = &b1

	case *Konfig:
		b1 := *bt
		a.KennungPtr = &b1

	case *Tai:
		b1 := *bt
		a.KennungPtr = &b1

	default:
		err = errors.Errorf("unsupported kennung: %T", b)
		return
	}

	return
}
