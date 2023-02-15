package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/values"
)

type Sigil int

const (
	SigilNone      = Sigil(iota)
	SigilSchwanzen = Sigil(1 << iota)
	SigilHistory
	SigilCwd
	SigilHidden

	SigilMax = SigilHidden
)

var (
	mapRuneToSigil = map[rune]Sigil{
		':': SigilNone,
		'@': SigilSchwanzen,
		'+': SigilHistory,
		'.': SigilCwd,
		'?': SigilHidden,
	}

	mapSigilToRune = map[Sigil]rune{
		SigilNone:      ':',
		SigilSchwanzen: '@',
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

func (a Sigil) GetGattung() schnittstellen.Gattung {
	return gattung.Unknown
}

func (a Sigil) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Sigil) Equals(b Sigil) bool {
	return a == b
}

func (a *Sigil) Reset() {
	*a = SigilNone
	return
}

func (a *Sigil) ResetWith(b Sigil) {
	*a = b
	return
}

func (a *Sigil) Add(b Sigil) {
	*a = *a | b
}

func (a *Sigil) Del(b Sigil) {
	*a = *a | ^b
}

func (a Sigil) Contains(b Sigil) bool {
	return a&b != 0
}

func (a Sigil) IncludesSchwanzen() bool {
	return a.Contains(SigilSchwanzen) || a.Contains(SigilHistory)
}

func (a Sigil) IncludesHistory() bool {
	return a.Contains(SigilHistory)
}

func (a Sigil) IncludesCwd() bool {
	return a.Contains(SigilCwd)
}

func (a Sigil) IncludesHidden() bool {
	return a.Contains(SigilHidden)
}

func (a Sigil) String() string {
	sb := strings.Builder{}
	errors.TodoP0("use sigil map")

	sb.WriteString(":")

	for s := SigilNone; s <= SigilMax; s++ {
		if a.Contains(s) {
			r := mapSigilToRune[s]
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func (i *Sigil) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	els := []rune(v)

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

func (i Sigil) GetSha() sha.Sha {
	return sha.FromString(i.String())
}
