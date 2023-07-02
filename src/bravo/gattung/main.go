package gattung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
)

type Gattung int

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

func (g Gattung) GetGattungString() string {
	return g.String()
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
	return strings.HasPrefix(v, p) || v == p
}

func (g *Gattung) Set(v string) (err error) {
	v = strings.TrimSpace(strings.ToLower(v))

	switch {
	case v == "akte":
		*g = Akte

	case hasPrefixOrEquals("typ", v):
		*g = Typ

	case v == "aktetyp":
		*g = Typ

	case hasPrefixOrEquals("etikett", v):
		*g = Etikett

	case hasPrefixOrEquals("zettel", v):
		*g = Zettel

	case v == "bezeichnung":
		*g = Bezeichnung

	case v == "hinweis":
		*g = Hinweis

	case v == "transaktion":
		*g = Transaktion

	case v == "konfig":
		*g = Konfig

	case v == "kennung":
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
