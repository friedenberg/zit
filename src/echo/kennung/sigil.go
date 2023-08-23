package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Sigil int

const (
	SigilSchwanzen = Sigil(iota)
	SigilHistory   = Sigil(1 << iota)
	SigilCwd
	SigilHidden

	SigilMax = SigilHidden
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
