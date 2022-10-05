package sha

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"path"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
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

func FromFormatString(f string, vs ...interface{}) Sha {
	return FromString(fmt.Sprintf(f, vs...))
}

func FromString(s string) Sha {
	hash := sha256.New()
	sr := strings.NewReader(s)

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return FromHash(hash)
}

func FromHash(h hash.Hash) (s Sha) {
	s = Sha{}
	s.Value = fmt.Sprintf("%x", h.Sum(nil))

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

func (s *Sha) SetParts(a, b string) (err error) {
	if err = s.Set(a + b); err != nil {
		err = errors.Wrap(err)
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

func (s Sha) Kopf() string {
	return s.String()[0:2]
}

func (s Sha) Schwanz() string {
	return s.String()[2:]
}

func (a Sha) Equals(b Sha) bool {
	return a.String() == b.String()
}

func (s Sha) Path(pc ...string) string {
	pc = append(pc, s.Kopf())
	pc = append(pc, s.Schwanz())

	return path.Join(pc...)
}
