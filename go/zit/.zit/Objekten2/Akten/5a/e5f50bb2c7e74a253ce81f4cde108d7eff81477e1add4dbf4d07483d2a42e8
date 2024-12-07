package genres

import (
	"fmt"
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
	None = Genre(iota)
	Blob
	Type
	_ // Bezeichnung
	Tag
	_ // Hinweis
	_ // Transaktion
	Zettel
	Config
	_ // Kennung
	InventoryList
	_ // AkteTyp
	Repo

	MaxGenre = Repo
)

const (
	unknown = byte(iota)
	blob    = byte(1 << iota)
	tipe
	tag
	zettel
	config
	inventory_list
	repo
)

// TODO switch to
// quiter.Slice[Genre](All())
func All() (out []Genre) {
	out = make([]Genre, 0, MaxGenre-1)

	for i := None + 1; i <= MaxGenre; i++ {
		out = append(out, Genre(i))
	}

	return
}

// TODO switch to
// quiter.Slice[Genre](TrueGenre())
func TrueGenre() (out []Genre) {
	out = make([]Genre, 0, MaxGenre-1)

	for i := None + 1; i <= MaxGenre; i++ {
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
		g = None
	}

	return
}

func (g Genre) GetGenre() interfaces.Genre {
	return g
}

func (g Genre) GetGenreBitInt() byte {
	switch g {
	default:
		panic(fmt.Sprintf("genre does not define bit int: %s", g))
	case InventoryList:
		return inventory_list
	case Blob:
		return blob
	case Zettel:
		return zettel
	case Tag:
		return tag
	case Repo:
		return repo
	case Type:
		return tipe
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
		err = MakeErrUnsupportedGenre(b)
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
	case Type, Tag, Zettel, Config, Repo, InventoryList, Blob:
		return true

	default:
		return false
	}
}

func (g Genre) GetGenreStringPlural(sv interfaces.StoreVersion) string {
	v := sv.GetInt()

	switch {
	case v <= 6:
		return g.getGenreStringPluralOld()

	default:
		return g.getGenreStringPluralNew()
	}
}

func (g Genre) getGenreStringPluralNew() string {
	switch g {
	case Blob:
		return "blobs"

	case Type:
		return "types"

	case Tag:
		return "tags"

	case Zettel:
		return "zettels"

	case InventoryList:
		return "inventory_lists"

	case Repo:
		return "repos"

	default:
		panic(fmt.Sprintf("undeclared plural for genre: %q", g))
	}
}

func (g Genre) getGenreStringPluralOld() string {
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

	case None:
		return "none"

	default:
		return fmt.Sprintf("Unknown(%#v)", g)
	}
}

func hasPrefixOrEquals(v, p string) bool {
	return strings.HasPrefix(v, p) || strings.EqualFold(v, p)
}

func (g *Genre) Set(v string) (err error) {
	v = strings.TrimSpace(v)

	switch {
	case strings.EqualFold(v, "blob"):
		fallthrough
	case strings.EqualFold(v, "akte"):
		*g = Blob

	case hasPrefixOrEquals("typ", v):
		*g = Type

	case strings.EqualFold(v, "aktetyp"):
		*g = Type

	case hasPrefixOrEquals("tag", v):
		fallthrough
	case hasPrefixOrEquals("etikett", v):
		*g = Tag

	case hasPrefixOrEquals("zettel", v):
		*g = Zettel

	case strings.EqualFold("config", v):
		fallthrough
	case strings.EqualFold("konfig", v):
		*g = Config

	case hasPrefixOrEquals("inventory_list", v):
		fallthrough
	case hasPrefixOrEquals("inventory-list", v):
		fallthrough
	case hasPrefixOrEquals("bestandsaufnahme", v):
		*g = InventoryList

	case hasPrefixOrEquals("repo", v):
		fallthrough
	case hasPrefixOrEquals("kasten", v):
		*g = Repo

	default:
		err = errors.Wrap(MakeErrUnrecognizedGenre(v))
		return
	}

	return
}

func (g *Genre) Reset() {
	*g = None
}

func (g *Genre) ReadFrom(r io.Reader) (n int64, err error) {
	*g = None

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
