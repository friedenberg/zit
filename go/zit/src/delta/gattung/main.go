package gattung

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

type Genre byte

// Do not change this order, various serialization formats rely on the
// underlying integer values.
const (
	Unknown = Genre(iota)
	Akte
	Typ
	Bezeichnung
	Etikett
	_ //Hinweis
	_ //Transaktion
	Zettel
	Konfig
	_ //Kennung
	Bestandsaufnahme
	AkteTyp
	Kasten

	MaxGattung = Kasten
)

const (
	unknown = byte(iota)
	akte    = byte(1 << iota)
	typ
	etikett
	zettel
	konfig
	kasten
)

func All() (out []Genre) {
	out = make([]Genre, 0, MaxGattung-1)

	for i := Unknown + 1; i <= MaxGattung; i++ {
		out = append(out, Genre(i))
	}

	return
}

func TrueGattung() (out []Genre) {
	out = make([]Genre, 0, MaxGattung-1)

	for i := Unknown + 1; i <= MaxGattung; i++ {
		g := Genre(i)

		if !g.IsTrueGattung() {
			continue
		}

		out = append(out, g)
	}

	return
}

func Must(g interfaces.GenreGetter) Genre {
	return g.GetGenre().(Genre)
}

func Make(g interfaces.Genre) Genre {
	return Must(g)
}

func MakeOrUnknown(v string) (g Genre) {
	if err := g.Set(v); err != nil {
		g = Unknown
	}

	return
}

func (g Genre) GetGenre() interfaces.Genre {
	return g
}

func (g Genre) GetGenreBitInt() byte {
	switch g {
	default:
		return unknown
	case Akte:
		return akte
	case Zettel:
		return zettel
	case Etikett:
		return etikett
	case Kasten:
		return kasten
	case Typ:
		return typ
	case Konfig:
		return konfig
	}
}

func (a Genre) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Genre) Equals(b Genre) bool {
	return a == b
}

func (a Genre) EqualsGenre(b interfaces.GenreGetter) bool {
	return a.GetGenreString() == b.GetGenre().GetGenreString()
}

func (a Genre) AssertGattung(b interfaces.GenreGetter) (err error) {
	if a.GetGenreString() != b.GetGenre().GetGenreString() {
		err = MakeErrUnsupportedGattung(b)
		return
	}

	return
}

func (g Genre) GetGenreString() string {
	return g.String()
}

func (g Genre) HasParents() bool {
	switch g {
	case Typ, Etikett, Kasten:
		return true

	default:
		return false
	}
}

func (g Genre) IsTrueGattung() bool {
	switch g {
	case Typ, Etikett, Zettel, Konfig, Kasten:
		return true

	default:
		return false
	}
}

func (g Genre) GetGenreStringPlural() string {
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

	case Bestandsaufnahme:
		return "Bestandsaufnahmen"

	case Kasten:
		return "Kisten"

	default:
		return g.String()
	}
}

func (g Genre) String() string {
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

	case Konfig:
		return "Konfig"

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

func (g *Genre) Set(v string) (err error) {
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

	case strings.EqualFold("konfig", v):
		*g = Konfig

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

func (g *Genre) Reset() {
	*g = Unknown
}

func (g *Genre) ReadFrom(r io.Reader) (n int64, err error) {
	*g = Unknown

	var b [1]byte

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n += int64(n1)

	if err != nil {
		return
	}

	*g = Genre(b[0])

	return
}

func (g *Genre) WriteTo(w io.Writer) (n int64, err error) {
	b := [1]byte{byte(*g)}

	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, b[:])
	n += int64(n1)

	if err != nil {
		return
	}

	return
}
