package kennung

import (
	"io"
	"math"
	"slices"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

var poolKennung3 interfaces.Pool[Kennung3, *Kennung3]

func init() {
	poolKennung3 = pool.MakePool(
		nil,
		func(k *Kennung3) {
			k.Reset()
		},
	)
}

func GetKennung3Pool() interfaces.Pool[Kennung3, *Kennung3] {
	return poolKennung3
}

type Kennung3 struct {
	g                   gattung.Gattung
	middle              byte // remove and replace with virtual
	kasten, left, right catgut.String
}

func MustKennung3(kp Kennung) (k *Kennung3) {
	k = &Kennung3{}
	err := k.SetWithKennung(kp)
	errors.PanicIfError(err)
	return
}

func (a *Kennung3) GetKasten() interfaces.KastenLike {
	return MustKasten(a.kasten.String())
}

func (a *Kennung3) IsVirtual() bool {
	switch a.g {
	case gattung.Zettel:
		return slices.Equal(a.left.Bytes(), []byte{'%'})

	case gattung.Etikett:
		return a.middle == '%' || slices.Equal(a.left.Bytes(), []byte{'%'})

	default:
		return false
	}
}

func (a *Kennung3) Equals(b *Kennung3) bool {
	if a.g != b.g {
		return false
	}

	if a.middle != b.middle {
		return false
	}

	if !a.left.Equals(&b.left) {
		return false
	}

	if !a.right.Equals(&b.right) {
		return false
	}

	if !a.kasten.Equals(&b.kasten) {
		return false
	}

	return true
}

func (k3 *Kennung3) WriteTo(w io.Writer) (n int64, err error) {
	if k3.Len() > math.MaxUint8 {
		err = errors.Errorf(
			"%q is greater than max uint8 (%d)",
			k3.String(),
			math.MaxUint8,
		)

		return
	}

	var n1 int64
	n1, err = k3.g.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	b := [2]uint8{uint8(k3.Len()), uint8(k3.left.Len())}

	var n2 int
	n2, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k3.left.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	bMid := [1]byte{k3.middle}

	n2, err = ohio.WriteAllOrDieTrying(w, bMid[:])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = k3.right.WriteTo(w)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *Kennung3) ReadFrom(r io.Reader) (n int64, err error) {
	var n1 int64
	n1, err = k3.g.ReadFrom(r)
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

	if _, err = k3.left.ReadNFrom(r, int(middlePos)); err != nil {
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

	k3.middle = bMiddle[0]

	if _, err = k3.right.ReadNFrom(r, int(contentLength-middlePos-1)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *Kennung3) SetGattung(g interfaces.GattungGetter) {
	if g == nil {
		k3.g = gattung.Unknown
	} else {
		k3.g = gattung.Must(g.GetGattung())
	}

	if k3.g == gattung.Zettel {
		k3.middle = '/'
	}
}

func (k3 *Kennung3) StringFromPtr() string {
	var sb strings.Builder

	switch k3.g {
	case gattung.Zettel:
		sb.Write(k3.left.Bytes())
		sb.WriteByte(k3.middle)
		sb.Write(k3.right.Bytes())

	case gattung.Typ:
		sb.Write(k3.right.Bytes())

	default:
		if k3.left.Len() > 0 {
			sb.Write(k3.left.Bytes())
		}

		if k3.middle != '\x00' {
			sb.WriteByte(k3.middle)
		}

		if k3.right.Len() > 0 {
			sb.Write(k3.right.Bytes())
		}
	}

	return sb.String()
}

func (k3 *Kennung3) IsEmpty() bool {
	if k3.g == gattung.Zettel {
		if k3.left.IsEmpty() && k3.right.IsEmpty() {
			return true
		}
	}

	return k3.left.Len() == 0 && k3.middle == 0 && k3.right.Len() == 0
}

func (k3 *Kennung3) Len() int {
	return k3.left.Len() + 1 + k3.right.Len()
}

func (k3 *Kennung3) KopfUndSchwanz() (kopf, schwanz string) {
	kopf = k3.left.String()
	schwanz = k3.right.String()

	return
}

func (k3 *Kennung3) LenKopfUndSchwanz() (int, int) {
	return k3.left.Len(), k3.right.Len()
}

func (k3 *Kennung3) String() string {
	return k3.StringFromPtr()
}

func (k3 *Kennung3) Reset() {
	k3.g = gattung.Unknown
	k3.left.Reset()
	k3.middle = 0
	k3.right.Reset()
}

func (k3 *Kennung3) PartsStrings() KennungParts {
	return KennungParts{
		Kasten: &k3.kasten,
		Left:   &k3.left,
		Middle: k3.middle,
		Right:  &k3.right,
	}
}

func (k3 *Kennung3) Parts() [3]string {
	var mid string

	if k3.middle != 0 {
		mid = string([]byte{k3.middle})
	}

	return [3]string{
		k3.left.String(),
		mid,
		k3.right.String(),
	}
}

func (k3 *Kennung3) GetGattung() interfaces.GattungLike {
	return k3.g
}

func MakeKennung3(
	v interfaces.StringerGattungGetter,
	ka Kasten,
) (k *Kennung3, err error) {
	k = &Kennung3{
		g: gattung.Unknown,
	}

	if err = k.kasten.Set(ka.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = k.SetWithGattung(v.String(), v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *Kennung3) Expand(
	a Abbr,
) (err error) {
	ex := a.ExpanderFor(k3.g)

	if ex == nil {
		return
	}

	v := k3.String()

	if v, err = ex(v); err != nil {
		err = nil
		return
	}

	if err = k3.SetWithGattung(v, k3.g); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (k3 *Kennung3) Abbreviate(
	a Abbr,
) (err error) {
	return
}

func (k3 *Kennung3) SetFromPath(
	path string,
	fe file_extensions.FileExtensions,
) (err error) {
	els := files.PathElements(path)
	ext := els[0]

	switch ext {
	case fe.Etikett:
		if err = k3.SetWithGattung(els[1], gattung.Etikett); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Typ:
		if err = k3.SetWithGattung(els[1], gattung.Typ); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Kasten:
		if err = k3.SetWithGattung(els[1], gattung.Kasten); err != nil {
			err = errors.Wrap(err)
			return
		}

	case fe.Zettel:
		if err = k3.SetWithGattung(els[2]+"/"+els[1], gattung.Zettel); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = ErrFDNotKennung
		return
	}

	return
}

func (h *Kennung3) SetWithKennung(
	k Kennung,
) (err error) {
	switch kt := k.(type) {
	case *Kennung3:
		if err = kt.left.CopyTo(&h.left); err != nil {
			err = errors.Wrap(err)
			return
		}

		h.middle = kt.middle

		if err = kt.right.CopyTo(&h.right); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		p := k.Parts()

		if err = h.left.Set(p[0]); err != nil {
			err = errors.Wrap(err)
			return
		}

		mid := []byte(p[1])

		if len(mid) >= 1 {
			h.middle = mid[0]
		}

		if err = h.right.Set(p[2]); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	h.SetGattung(k)

	return
}

func (h *Kennung3) SetWithGattung(
	v string,
	g interfaces.GattungGetter,
) (err error) {
	h.g = gattung.Make(g.GetGattung())

	if err = h.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (h *Kennung3) TodoSetBytes(v *catgut.String) (err error) {
	return h.Set(v.String())
}

func (h *Kennung3) Set(v string) (err error) {
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
		err = errors.Wrapf(err, "String: %q", v)
		return
	}

	if err = h.SetWithKennung(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a *Kennung3) ResetWith(b *Kennung3) {
	a.g = b.g
	b.left.CopyTo(&a.left)
	b.right.CopyTo(&a.right)
	a.middle = b.middle
}

func (a *Kennung3) ResetWithKennung(b Kennung) (err error) {
	return a.SetWithKennung(b)
}

func (t *Kennung3) MarshalText() (text []byte, err error) {
	text = []byte(FormattedString(t))
	return
}

func (t *Kennung3) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t *Kennung3) MarshalBinary() (text []byte, err error) {
	// if t.g == gattung.Unknown {
	// 	err = errors.Wrapf(gattung.ErrEmptyKennung{}, "Kennung: %s", t)
	// 	return
	// }

	text = []byte(FormattedString(t))

	return
}

func (t *Kennung3) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
