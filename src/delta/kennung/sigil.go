package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type Sigil int

const (
	SigilNone = Sigil(iota)
	SigilAll  = Sigil(1 << iota)
	SigilHistory
	SigilCwd
)

var (
	sigilMap = map[rune]Sigil{
		'!': SigilAll,
		'+': SigilHistory,
		'@': SigilCwd,
	}
)

func MakeSigil(v Sigil) (s *Sigil) {
	s1 := Sigil(v)
	s = &s1
	return
}

func (a Sigil) GetGattung() schnittstellen.Gattung {
	return gattung.Unknown
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
	*a = *a |^ b
}

func (a Sigil) Contains(b Sigil) bool {
	return a&b != 0
}

func (a Sigil) String() string {
	sb := strings.Builder{}

	if a.Contains(SigilAll) {
		sb.WriteString("!")
	}

	if a.Contains(SigilHistory) {
		sb.WriteString("+")
	}

	if a.Contains(SigilCwd) {
		sb.WriteString("@")
	}

	return sb.String()
}

func (i *Sigil) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	els := []rune(v)

	if len(els) > 3 {
		err = errors.Errorf("not a sigil, too long")
		return
	}

	for _, v1 := range els {
		switch v1 {
		case '!', '+', '@':
			i.Add(sigilMap[v1])

		default:
			err = errors.Errorf("not a sigil")
			return
		}
	}

	return
}

func (i Sigil) GetSha() sha.Sha {
	return sha.FromString(i.String())
}
