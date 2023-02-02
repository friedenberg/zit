package gattung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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

	MaxGattung = AkteTyp
)

func All() (out []Gattung) {
	out = make([]Gattung, 0, MaxGattung-1)

	for i := Unknown + 1; i <= MaxGattung; i++ {
		out = append(out, Gattung(i))
	}

	return
}

func Make(g schnittstellen.Gattung) Gattung {
	return g.(Gattung)
}

func (g Gattung) GetGattung() schnittstellen.Gattung {
	return g
}

func (a Gattung) Equals(b1 any) bool {
	var b Gattung
	ok := false

	if b, ok = b1.(Gattung); !ok {
		return false
	}

	return a.GetGattungString() == b.GetGattungString()
}

func (g Gattung) GetGattungString() string {
	return g.String()
}

func (g Gattung) String() string {
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

	default:
		return "Unknown"
	}
}

func (g *Gattung) Set(v string) (err error) {
	v1 := strings.ToLower(v)

	switch v1 {
	case "akte":
		*g = Akte

	case "typ":
		*g = Typ

	case "aktetyp":
		*g = Typ

	case "etikett":
		*g = Etikett

	case "zettel":
		*g = Zettel

	case "bezeichnung":
		*g = Bezeichnung

	case "hinweis":
		*g = Hinweis

	case "transaktion":
		*g = Transaktion

	case "konfig":
		*g = Konfig

	case "kennung":
		*g = Kennung

	case "bestandsaufnahme":
		*g = Bestandsaufnahme

	default:
		err = errors.Wrap(ErrUnrecognizedGattung{string: v1})
		return
	}

	return
}
