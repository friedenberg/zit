package kennung

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

var poolKennung2 schnittstellen.Pool[Kennung2, *Kennung2]

func init() {
	poolKennung2 = pool.MakePool[Kennung2, *Kennung2](
		nil,
		func(k *Kennung2) {
			k.Reset()
		},
	)
}

func GetKennungPool() schnittstellen.Pool[Kennung2, *Kennung2] {
	return poolKennung2
}

type Kennung2 struct {
	g                   gattung.Gattung
	left, middle, right catgut.String
}

func MustKennung2(kp Kennung) (k *Kennung2) {
	k = &Kennung2{}
	err := k.SetWithKennung(kp)
	errors.PanicIfError(err)
	return
}

func (k2 *Kennung2) SetGattung(g schnittstellen.GattungGetter) {
	if g == nil {
		k2.g = gattung.Unknown
	} else {
		k2.g = gattung.Must(g.GetGattung())
	}
}

func (k2 *Kennung2) StringFromPtr() string {
	var sb strings.Builder

	switch k2.g {
	case gattung.Zettel, gattung.Bestandsaufnahme, gattung.Kasten:
		sb.Write(k2.left.Bytes())
		sb.Write(k2.middle.Bytes())
		sb.Write(k2.right.Bytes())

	case gattung.Etikett, gattung.Typ, gattung.Konfig:
		sb.Write(k2.right.Bytes())

	default:
		sb.WriteString("unknown")
	}

	return sb.String()
}

func (k2 *Kennung2) IsEmpty() bool {
	return k2.left.Len() == 0 && k2.middle.Len() == 0 && k2.right.Len() == 0
}

func (k2 *Kennung2) Len() int {
	return k2.left.Len() + k2.middle.Len() + k2.right.Len()
}

func (k2 *Kennung2) KopfUndSchwanz() (kopf, schwanz string) {
	kopf = k2.left.String()
	schwanz = k2.right.String()

	return
}

func (k2 *Kennung2) LenKopfUndSchwanz() (int, int) {
	return k2.left.Len(), k2.right.Len()
}

func (src *Kennung2) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int64

	n1, err = src.left.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = src.middle.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = src.right.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k2 *Kennung2) String() string {
	return k2.StringFromPtr()
}

func (k2 *Kennung2) Reset() {
	k2.g = gattung.Unknown
	k2.left.Reset()
	k2.middle.Reset()
	k2.right.Reset()
}

func (k2 *Kennung2) PartsStrings() [3]*catgut.String {
	return [3]*catgut.String{
		&k2.left,
		&k2.middle,
		&k2.right,
	}
}

func (k2 *Kennung2) Parts() [3]string {
	return [3]string{
		k2.left.String(),
		k2.middle.String(),
		k2.right.String(),
	}
}

func (k2 *Kennung2) GetGattung() schnittstellen.GattungLike {
	return k2.g
}

func MakeKennung2(v string) (KennungPtr, error) {
	k := &Kennung2{
		g: gattung.Unknown,
	}

	return k, k.Set(v)
}

func (h *Kennung2) SetWithKennung(
	k Kennung,
) (err error) {
	switch kt := k.(type) {
	case *Kennung2:
		if err = kt.left.CopyTo(&h.left); err != nil {
			return
		}

		if err = kt.middle.CopyTo(&h.middle); err != nil {
			return
		}

		if err = kt.right.CopyTo(&h.right); err != nil {
			return
		}

	default:
		p := k.Parts()

		if err = h.left.Set(p[0]); err != nil {
			return
		}

		if err = h.middle.Set(p[1]); err != nil {
			return
		}

		if err = h.right.Set(p[2]); err != nil {
			return
		}
	}

	h.SetGattung(k)

	return
}

func (h *Kennung2) SetWithGattung(
	v string,
	g schnittstellen.GattungGetter,
) (err error) {
	h.g = gattung.Make(g.GetGattung())

	return h.Set(v)
}

func (h *Kennung2) TodoSetBytes(v *catgut.String) (err error) {
	return h.Set(v.String())
}

func (h *Kennung2) Set(v string) (err error) {
	var k Kennung

	switch h.g {
	case gattung.Unknown:
		k, err = Make(v)

	case gattung.Zettel:
		var h Hinweis
		err = h.Set(v)
		k = h

	case gattung.Etikett:
		var h Etikett
		err = h.Set(v)
		k = h

	case gattung.Typ:
		var h Typ
		err = h.Set(v)
		k = h

	case gattung.Kasten:
		var h Kasten
		err = h.Set(v)
		k = h

	case gattung.Konfig:
		var h Konfig
		err = h.Set(v)
		k = h

	case gattung.Bestandsaufnahme:
		var h Tai
		err = h.Set(v)
		k = h

	default:
		err = gattung.MakeErrUnrecognizedGattung(h.g.GetGattungString())
	}

	if err != nil {
		return
	}

	return h.SetWithKennung(k)
}

func (a *Kennung2) ResetWithKennung(b Kennung) (err error) {
	return a.SetWithKennung(b)
}

func (t *Kennung2) MarshalText() (text []byte, err error) {
	text = []byte(FormattedString(t))
	return
}

func (t *Kennung2) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *Kennung2) MarshalBinary() (text []byte, err error) {
	// if t.g == gattung.Unknown {
	// 	err = errors.Wrapf(gattung.ErrEmptyKennung{}, "Kennung: %s", t)
	// 	return
	// }

	text = []byte(FormattedString(t))

	return
}

func (t *Kennung2) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
