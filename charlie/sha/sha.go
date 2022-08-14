package sha

import (
	"encoding/hex"
	"fmt"
	"hash"
	"path"
	"strings"

	"github.com/friedenberg/zit/bravo/errors"
)

const (
	ShaNull = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

type PathComponents interface {
	PathComponents() []string
}

type Sha struct {
	Value string
}

func MakeSha(v string) (s Sha, err error) {
	s = Sha{}
	err = s.Set(v)

	return
}

func FromHash(h hash.Hash) (s Sha) {
	s = Sha{}
	s.SetFromHash(h)

	return
}

func (s Sha) String() string {
	if s.Value == "" {
		return ShaNull
	} else {
		return s.Value
	}
}

func (s Sha) Sha() Sha {
	return s
}

func (s *Sha) SetFromHash(h hash.Hash) {
	s.Value = fmt.Sprintf("%x", h.Sum(nil))
}

func (s *Sha) SetParts(a, b string) (err error) {
	if err = s.Set(a + b); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s *Sha) Set(v string) (err error) {
	v1 := strings.TrimSpace(v)

	if _, err = hex.DecodeString(v1); err != nil {
		err = errors.Errorf("%q: %s", v1, err)
		return
	}

	s.Value = v1

	return
}

func (s Sha) IsNull() bool {
	if s.Value == "" {
		return true
	}

	if s.Value == ShaNull {
		return true
	}

	return false
}

func (s Sha) Head() string {
	return s.String()[0:2]
}

func (s Sha) Tail() string {
	return s.String()[2:]
}

func (a Sha) Equals(b Sha) bool {
	return a.String() == b.String()
}

func (s Sha) Path(pc ...string) string {
	pc = append(pc, s.Head())
	pc = append(pc, s.Tail())

	return path.Join(pc...)
}
