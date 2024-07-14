package genres

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
	Blob
	Type
	_ //Bezeichnung
	Tag
	_ //Hinweis
	_ //Transaktion
	Zettel
	Config
	_ //Kennung
	InventoryList
	_ // AkteTyp
	Repo

	MaxGenre = Repo
)

const (
	unknown = byte(iota)
	akte    = byte(1 << iota)
	typ
	etikett
	zettel
	config
	kasten
)

func All() (out []Genre) {
	out = make([]Genre, 0, MaxGenre-1)

	for i := Unknown + 1; i <= MaxGenre; i++ {
		out = append(out, Genre(i))
	}

	return
}

func TrueGenre() (out []Genre) {
	out = make([]Genre, 0, MaxGenre-1)

	for i := Unknown + 1; i <= MaxGenre; i++ {
		g := Genre(i)

		if !g.IsTrueGenre() {
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
	case Blob:
		return akte
	case Zettel:
		return zettel
	case Tag:
		return etikett
	case Repo:
		return kasten
	case Type:
		return typ
	case Config:
		return config
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

func (a Genre) AssertGenre(b interfaces.GenreGetter) (err error) {
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
	case Type, Tag, Repo:
		return true

	default:
		return false
	}
}

func (g Genre) IsTrueGenre() bool {
	switch g {
	case Type, Tag, Zettel, Config, Repo:
		return true

	default:
		return false
	}
}

func (g Genre) GetGenreStringPlural() string {
	switch g {
	case Blob:
		return "Akten"

	case Type:
		return "Typen"

	case Tag:
		return "Etiketten"

	case Zettel:
		return "Zettelen"

	case InventoryList:
		return "Bestandsaufnahmen"

	case Repo:
		return "Kisten"

	default:
		return g.String()
	}
}

func (g Genre) String() string {
	switch g {
	case Blob:
		return "Akte"

	case Type:
		return "Typ"

	case Tag:
		return "Etikett"

	case Zettel:
		return "Zettel"

	case Config:
		return "Konfig"

	case InventoryList:
		return "Bestandsaufnahme"

	case Repo:
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
		*g = Blob

	case hasPrefixOrEquals("typ", v):
		*g = Type

	case strings.EqualFold(v, "aktetyp"):
		*g = Type

	case hasPrefixOrEquals("etikett", v):
		*g = Tag

	case hasPrefixOrEquals("zettel", v):
		*g = Zettel

	case strings.EqualFold("konfig", v):
		*g = Config

	case hasPrefixOrEquals("bestandsaufnahme", v):
		*g = InventoryList

	case hasPrefixOrEquals("kasten", v):
		*g = Repo

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
