package gattung

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/delta/ohio"
)

type Gattung byte

// Do not change this order, various serialization formats rely on the
// underlying integer values.
const (
	Unknown = Gattung(iota)
	Akte
	Typ
	Bezeichnung
	Etikett
	Hinweis
	Transaktion
	Zettel
	Konfig
	Kennung
	Bestandsaufnahme
	AkteTyp
	Kasten

	MaxGattung = Kasten
)

func All() (out []Gattung) {
	out = make([]Gattung, 0, MaxGattung-1)

	for i := Unknown + 1; i <= MaxGattung; i++ {
		out = append(out, Gattung(i))
	}

	return
}

func TrueGattung() (out []Gattung) {
	out = make([]Gattung, 0, MaxGattung-1)

	for i := Unknown + 1; i <= MaxGattung; i++ {
		g := Gattung(i)

		if !g.IsTrueGattung() {
			continue
		}

		out = append(out, g)
	}

	return
}

func Must(g schnittstellen.GattungGetter) Gattung {
	return g.GetGattung().(Gattung)
}

func Make(g schnittstellen.GattungLike) Gattung {
	return Must(g)
}

func MakeOrUnknown(v string) (g Gattung) {
	if err := g.Set(v); err != nil {
		g = Unknown
	}

	return
}

func (g Gattung) GetGattung() schnittstellen.GattungLike {
	return g
}

func (a Gattung) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Gattung) Equals(b Gattung) bool {
	return a == b
}

func (a Gattung) EqualsGattung(b schnittstellen.GattungGetter) bool {
	return a.GetGattungString() == b.GetGattung().GetGattungString()
}

func (a Gattung) AssertGattung(b schnittstellen.GattungGetter) (err error) {
	if a.GetGattungString() != b.GetGattung().GetGattungString() {
		err = MakeErrUnsupportedGattung(b)
		return
	}

	return
}

func (g Gattung) GetGattungString() string {
	return g.String()
}

func (g Gattung) HasParents() bool {
	switch g {
	case Typ, Etikett, Kasten:
		return true

	default:
		return false
	}
}

func (g Gattung) IsTrueGattung() bool {
	switch g {
	case Typ, Etikett, Zettel, Konfig, Kasten:
		return true

	default:
		return false
	}
}

func (g Gattung) GetGattungStringPlural() string {
	switch g {
	case Akte:
		return "Akten"

	case Typ:
		return "Typen"

	case Etikett:
		return "Etiketten"

	case Zettel:
		return "Zettelen"

	case Bezeichnung:
		return "Bezeichnungen"

	case Hinweis:
		return "Hinweisen"

	case Kennung:
		return "Kennungen"

	case Bestandsaufnahme:
		return "Bestandsaufnahmen"

	case Kasten:
		return "Kisten"

	default:
		return g.String()
	}
}

func (g Gattung) String() string {
	errors.TodoP1(
		"move Bezeichnung, AkteTyp, Kennung, Transaktion, to another place",
	)
	switch g {
	case Akte:
		return "Akte"

	case Typ:
		return "Typ"

	case AkteTyp:
		return "AkteTyp"

	case Etikett:
		return "Etikett"

	case Zettel:
		return "Zettel"

	case Bezeichnung:
		return "Bezeichnung"

	case Hinweis:
		return "Hinweis"

	case Transaktion:
		return "Transaktion"

	case Konfig:
		return "Konfig"

	case Kennung:
		return "Kennung"

	case Bestandsaufnahme:
		return "Bestandsaufnahme"

	case Kasten:
		return "Kasten"

	default:
		return "Unknown"
	}
}

func hasPrefixOrEquals(v, p string) bool {
	return strings.HasPrefix(v, p) || strings.EqualFold(v, p)
}

func (g *Gattung) Set(v string) (err error) {
	v = strings.TrimSpace(v)

	switch {
	case strings.EqualFold(v, "akte"):
		*g = Akte

	case hasPrefixOrEquals("typ", v):
		*g = Typ

	case strings.EqualFold(v, "aktetyp"):
		*g = Typ

	case hasPrefixOrEquals("etikett", v):
		*g = Etikett

	case hasPrefixOrEquals("zettel", v):
		*g = Zettel

	case strings.EqualFold(v, "bezeichnung"):
		*g = Bezeichnung

	case strings.EqualFold("hinweis", v):
		*g = Hinweis

	case strings.EqualFold("transaktion", v):
		*g = Transaktion

	case strings.EqualFold("konfig", v):
		*g = Konfig

	case strings.EqualFold("kennung", v):
		*g = Kennung

	case hasPrefixOrEquals("bestandsaufnahme", v):
		*g = Bestandsaufnahme

	case hasPrefixOrEquals("kasten", v):
		*g = Kasten

	default:
		err = errors.Wrap(MakeErrUnrecognizedGattung(v))
		return
	}

	return
}

func (g *Gattung) Reset() {
	*g = Unknown
}

func (g *Gattung) ReadFrom(r io.Reader) (n int64, err error) {
	*g = Unknown

	var b [1]byte

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n1)

	if err != nil {
		return
	}

	*g = Gattung(b[0])

	return
}

func (g *Gattung) WriteTo(w io.Writer) (n int64, err error) {
	b := [1]byte{byte(*g)}

	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n1)

	if err != nil {
		return
	}

	return
}
