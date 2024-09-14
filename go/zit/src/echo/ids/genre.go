package ids

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/token_types"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
)

type Genre byte

func MakeGenreAll() Genre {
	return MakeGenre(genres.TrueGenre()...)
}

func MakeGenre(vs ...genres.Genre) (s Genre) {
	s.Add(vs...)
	return
}

func (a Genre) IsEmpty() bool {
	return a == 0
}

func (a Genre) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Genre) Equals(b Genre) bool {
	return a == b
}

func (a *Genre) Reset() {
	*a = 0
}

func (a *Genre) ResetWith(b Genre) {
	*a = b
}

func (a *Genre) Add(bs ...genres.Genre) {
	for _, b := range bs {
		*a |= Genre(b.GetGenre().GetGenreBitInt())
	}
}

func (a *Genre) Del(b interfaces.GenreGetter) {
	*a &= ^Genre(b.GetGenre().GetGenreBitInt())
}

func (a Genre) Contains(b interfaces.GenreGetter) bool {
	bg := Genre(b.GetGenre().GetGenreBitInt())
	return byte(a&bg) == byte(bg)
}

func (a Genre) ContainsOneOf(b interfaces.GenreGetter) bool {
	bg := Genre(b.GetGenre().GetGenreBitInt())
	return a&bg != 0
}

func (a Genre) Slice() []genres.Genre {
	tg := genres.TrueGenre()
	out := make([]genres.Genre, 0, len(tg))

	for _, g := range tg {
		if !a.ContainsOneOf(g) {
			continue
		}

		out = append(out, g)
	}

	return out
}

func (a Genre) String() string {
	sb := strings.Builder{}

	first := true

	for _, g := range genres.TrueGenre() {
		if !a.ContainsOneOf(g) {
			continue
		}

		if !first {
			sb.WriteRune(',')
		}

		sb.WriteString(g.GetGenreString())
		first = false
	}

	return sb.String()
}

func (i *Genre) AddString(v string) (err error) {
	var g genres.Genre

	if err = g.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.Add(g)

	return
}

func (gs *Genre) Set(v string) (err error) {
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

func (g *Genre) ReadFromTokenScanner(
	ts *query_spec.TokenScanner,
) (err error) {
	for ts.Scan() {
		token, tokenType := ts.GetTokenAndType()

		switch tokenType {
		case token_types.TypeOperator:
			el := token.String()

			if el == " " {
				ts.Unscan()
				return
			}

			if el == "," {
				continue
			}

		case token_types.TypeIdentifier:
			if err = g.AddString(token.String()); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if err = ts.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i Genre) GetSha() *sha.Sha {
	return sha.FromString(i.String())
}

func (i Genre) Byte() byte {
	return byte(i)
}

func (i Genre) ReadByte() (byte, error) {
	return byte(i), nil
}

func (i *Genre) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n = int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	*i = Genre(b[0])

	return
}

func (i *Genre) WriteTo(w io.Writer) (n int64, err error) {
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
