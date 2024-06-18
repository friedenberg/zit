package kennung

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type Gattung byte

func MakeGattungAll() Gattung {
	return MakeGattung(gattung.TrueGattung()...)
}

func MakeGattung(vs ...gattung.Gattung) (s Gattung) {
	s.Add(vs...)
	return
}

func (a Gattung) IsEmpty() bool {
	return a == 0
}

func (a Gattung) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Gattung) Equals(b Gattung) bool {
	return a == b
}

func (a *Gattung) Reset() {
	*a = 0
}

func (a *Gattung) ResetWith(b Gattung) {
	*a = b
}

func (a *Gattung) Add(bs ...gattung.Gattung) {
	for _, b := range bs {
		*a |= Gattung(b.GetGattung().GetGattungBitInt())
	}
}

func (a *Gattung) Del(b schnittstellen.GattungGetter) {
	*a &= ^Gattung(b.GetGattung().GetGattungBitInt())
}

func (a Gattung) Contains(b schnittstellen.GattungGetter) bool {
	bg := Gattung(b.GetGattung().GetGattungBitInt())
	return byte(a&bg) == byte(bg)
}

func (a Gattung) ContainsOneOf(b schnittstellen.GattungGetter) bool {
	bg := Gattung(b.GetGattung().GetGattungBitInt())
	return a&bg != 0
}

func (a Gattung) Slice() []gattung.Gattung {
	tg := gattung.TrueGattung()
	out := make([]gattung.Gattung, 0, len(tg))

	for _, g := range tg {
		if !a.ContainsOneOf(g) {
			continue
		}

		out = append(out, g)
	}

	return out
}

func (a Gattung) String() string {
	sb := strings.Builder{}

	first := true

	for _, g := range gattung.TrueGattung() {
		if !a.ContainsOneOf(g) {
			continue
		}

		if !first {
			sb.WriteRune(',')
		}

		sb.WriteString(g.GetGattungString())
		first = false
	}

	return sb.String()
}

func (i *Gattung) AddString(v string) (err error) {
	var g gattung.Gattung

	if err = g.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.Add(g)

	return
}

func (gs *Gattung) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	for _, g := range strings.Split(v, ",") {
		if err = gs.AddString(g); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (g *Gattung) SetTokens(
	tokens ...string,
) (remainingTokens []string, err error) {
	for i, el := range tokens {
		if el == " " {
			remainingTokens = tokens[i:]
			return
		}

		if el == "," {
			continue
		}

		if err = g.AddString(el); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i Gattung) GetSha() *sha.Sha {
	return sha.FromString(i.String())
}

func (i Gattung) Byte() byte {
	return byte(i)
}

func (i Gattung) ReadByte() (byte, error) {
	return byte(i), nil
}

func (i *Gattung) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n = int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	*i = Gattung(b[0])

	return
}

func (i *Gattung) WriteTo(w io.Writer) (n int64, err error) {
	var b byte

	if b, err = i.ReadByte(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, []byte{b})
	n = int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
