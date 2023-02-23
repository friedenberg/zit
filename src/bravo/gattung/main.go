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

func Must(g schnittstellen.Gattung) Gattung {
	return g.(Gattung)
}

func Make(g schnittstellen.Gattung) Gattung {
	return g.(Gattung)
}

func (g Gattung) GetGattung() schnittstellen.Gattung {
	return g
}

func (a Gattung) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Gattung) Equals(b Gattung) bool {
	return a == b
}

func (g Gattung) GetGattungString() string {
	return g.String()
}

func (g Gattung) String() string {
	errors.TodoP0("determine of some of these gattung should be goodbyed")
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

func (g *Gattung) Set(v string) (err error) {
	v = strings.TrimSpace(strings.ToLower(v))

	switch {
	case v == "akte":
		*g = Akte

	case strings.HasPrefix("typ", v):
		*g = Typ

	case v == "aktetyp":
		*g = Typ

	case strings.HasPrefix("etikett", v):
		*g = Etikett

	case strings.HasPrefix("zettel", v):
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

	case strings.HasPrefix("bestandsaufnahme", v):
		*g = Bestandsaufnahme

	case strings.HasPrefix("kasten", v):
		*g = Kasten

	default:
		err = errors.Wrap(ErrUnrecognizedGattung{string: v})
		return
	}

	return
}
