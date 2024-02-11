package kennung

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/values"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/delta/ohio"
)

type Sigil byte

const (
	SigilUnknown   = Sigil(iota)
	SigilSchwanzen = Sigil(1 << iota)
	SigilHistory
	SigilCwd
	SigilHidden

	SigilMax
	SigilAll = Sigil(^byte(0))
)

var (
	mapRuneToSigil = map[rune]Sigil{
		':': SigilSchwanzen,
		'+': SigilHistory,
		'.': SigilCwd,
		'?': SigilHidden,
	}

	mapSigilToRune = map[Sigil]rune{
		SigilSchwanzen: ':',
		SigilHistory:   '+',
		SigilCwd:       '.',
		SigilHidden:    '?',
	}
)

func SigilFieldFunc(c rune) (ok bool) {
	_, ok = mapRuneToSigil[c]
	return
}

func MakeSigil(v Sigil) (s *Sigil) {
	s1 := Sigil(v)
	s = &s1
	return
}

func (a Sigil) GetGattung() schnittstellen.GattungLike {
	return gattung.Unknown
}

func (a Sigil) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Sigil) Equals(b Sigil) bool {
	return a == b
}

func (a *Sigil) Reset() {
	*a = SigilSchwanzen
}

func (a *Sigil) ResetWith(b Sigil) {
	*a = b
}

func (a *Sigil) Add(b Sigil) {
	*a |= b
}

func (a *Sigil) Del(b Sigil) {
	*a &= ^b
}

func (a Sigil) Contains(b Sigil) bool {
	return a&b != 0
}

func (a Sigil) GetSigil() schnittstellen.Sigil {
	return a
}

func (a Sigil) IncludesSchwanzen() bool {
	return a.Contains(SigilSchwanzen) || a.Contains(SigilHistory) || a == 0
}

func (a Sigil) IncludesHistory() bool {
	return a.Contains(SigilHistory)
}

func (a Sigil) IncludesCwd() bool {
	return a.Contains(SigilCwd)
}

func (a Sigil) IncludesHidden() bool {
	return a.Contains(SigilHidden) || a.Contains(SigilCwd)
}

func (a Sigil) String() string {
	sb := strings.Builder{}

	for s := SigilSchwanzen; s <= SigilMax; s++ {
		if a.Contains(s) {
			r, ok := mapSigilToRune[s]

			if !ok {
				continue
			}

			sb.WriteRune(r)
		}
	}

	if sb.Len() == 0 {
		sb.WriteRune(mapSigilToRune[SigilSchwanzen])
	}

	return sb.String()
}

func (i *Sigil) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	els := v

	for _, v1 := range els {
		if _, ok := mapRuneToSigil[v1]; ok {
			i.Add(mapRuneToSigil[v1])
		} else {
			err = errors.Wrap(errInvalidSigil(v))
			return
		}
	}

	return
}

func (i Sigil) GetSha() *sha.Sha {
	return sha.FromString(i.String())
}

func (i Sigil) Byte() byte {
	return byte(i)
}

func (i Sigil) ReadByte() (byte, error) {
	return byte(i), nil
}

func (i *Sigil) ReadFrom(r io.Reader) (n int64, err error) {
	var b [1]byte

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, b[:])
	n = int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	*i = Sigil(b[0])

	return
}

func (i *Sigil) WriteTo(w io.Writer) (n int64, err error) {
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
