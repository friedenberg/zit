package kennung

import (
	"io"
	"math"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/src/charlie/catgut"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/delta/ohio"
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
	g           gattung.Gattung
	middle      byte
	left, right catgut.String
}

func MustKennung2(kp Kennung) (k *Kennung2) {
	k = &Kennung2{}
	err := k.SetWithKennung(kp)
	errors.PanicIfError(err)
	return
}

func (k2 *Kennung2) WriteTo(w io.Writer) (n int64, err error) {
	if k2.Len() > math.MaxUint8 {
		err = errors.Errorf(
			"%q is greater than max uint8 (%d)",
			k2.String(),
			math.MaxUint8,
		)

		return
	}

	var n1 int64
	n1, err = k2.g.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b := [2]uint8{uint8(k2.Len()), uint8(k2.left.Len())}

	var n2 int
	n2, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k2.left.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	bMid := [1]byte{k2.middle}

	n2, err = ohio.WriteAllOrDieTrying(w, bMid[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k2.right.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k2 *Kennung2) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = k2.g.ReadFrom(r)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var b [2]uint8

	var n2 int
	n2, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	contentLength := b[0]
	middlePos := b[1]

	if middlePos > contentLength-1 {
		err = errors.Errorf(
			"middle position %d is greater than last index: %d",
			middlePos,
			contentLength,
		)
		return
	}

	if _, err = k2.left.ReadNFrom(r, int(middlePos)); err != nil {
		err = errors.Wrap(err)
		return
	}

	var bMiddle [1]uint8

	n2, err = ohio.ReadAllOrDieTrying(r, bMiddle[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	k2.middle = bMiddle[0]

	if _, err = k2.right.ReadNFrom(r, int(contentLength-middlePos-1)); err != nil {
		err = errors.Wrap(err)
		return
	}

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
	case gattung.Zettel:
		sb.Write(k2.left.Bytes())
		sb.WriteByte(k2.middle)
		sb.Write(k2.right.Bytes())

	case gattung.Typ:
		sb.Write(k2.right.Bytes())

	default:
		if k2.left.Len() > 0 {
			sb.Write(k2.left.Bytes())
		}

		if k2.middle != '\x00' {
			sb.WriteByte(k2.middle)
		}

		if k2.right.Len() > 0 {
			sb.Write(k2.right.Bytes())
		}
	}

	return sb.String()
}

func (k2 *Kennung2) IsEmpty() bool {
	return k2.left.Len() == 0 && k2.middle == 0 && k2.right.Len() == 0
}

func (k2 *Kennung2) Len() int {
	return k2.left.Len() + 1 + k2.right.Len()
}

func (k2 *Kennung2) KopfUndSchwanz() (kopf, schwanz string) {
	kopf = k2.left.String()
	schwanz = k2.right.String()

	return
}

func (k2 *Kennung2) LenKopfUndSchwanz() (int, int) {
	return k2.left.Len(), k2.right.Len()
}

func (k2 *Kennung2) String() string {
	return k2.StringFromPtr()
}

func (k2 *Kennung2) Reset() {
	k2.g = gattung.Unknown
	k2.left.Reset()
	k2.middle = 0
	k2.right.Reset()
}

type KennungParts struct {
	Middle      byte
	Left, Right *catgut.String
}

func (k2 *Kennung2) PartsStrings() KennungParts {
	return KennungParts{
		Left:   &k2.left,
		Middle: k2.middle,
		Right:  &k2.right,
	}
}

func (k2 *Kennung2) Parts() [3]string {
	var mid string

	if k2.middle != 0 {
		mid = string([]byte{k2.middle})
	}

	return [3]string{
		k2.left.String(),
		mid,
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

		h.middle = kt.middle

		if err = kt.right.CopyTo(&h.right); err != nil {
			return
		}

	default:
		p := k.Parts()

		if err = h.left.Set(p[0]); err != nil {
			return
		}

		mid := []byte(p[1])

		if len(mid) >= 1 {
			h.middle = mid[0]
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

	if err = h.SetWithKennung(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
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
