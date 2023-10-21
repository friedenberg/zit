package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type Kennung2 struct {
	g     gattung.Gattung
	parts [3]string
}

func MustKennung2(kp KennungPtr) (k Kennung2) {
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

func (k2 Kennung2) String() string {
	var sb strings.Builder

	switch k2.g {
	case gattung.Zettel, gattung.Bestandsaufnahme:
		sb.WriteString(k2.parts[0])
		sb.WriteString(k2.parts[1])
		sb.WriteString(k2.parts[2])

	case gattung.Etikett, gattung.Typ, gattung.Kasten, gattung.Konfig:
		sb.WriteString(k2.parts[2])

	default:
		sb.WriteString("unknown")
	}

	return sb.String()
}

func (k2 *Kennung2) Reset() {
	k2.g = gattung.Unknown
	k2.parts[0] = ""
	k2.parts[1] = ""
	k2.parts[2] = ""
}

func (k2 Kennung2) Parts() [3]string {
	return k2.parts
}

func (k2 Kennung2) GetGattung() schnittstellen.GattungLike {
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
	h.parts = k.Parts()
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

func (a *Kennung2) ResetWithKennungPtr(b KennungPtr) (err error) {
	return a.SetWithKennung(b)
}

func (t Kennung2) MarshalText() (text []byte, err error) {
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

func (t Kennung2) MarshalBinary() (text []byte, err error) {
	if t.g == gattung.Unknown {
		err = gattung.ErrEmptyKennung{}
		return
	}

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
