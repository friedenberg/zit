package sha

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
)

const (
	ShaNull = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

type PathComponents interface {
	PathComponents() []string
}

type ShaLike = schnittstellen.ShaGetter

type Sha struct {
	value string
}

func MakeHashWriter() (h hash.Hash) {
	h = sha256.New()
	return
}

func Make(s schnittstellen.Sha) Sha {
	return s.(Sha)
}

func Must(v string) (s Sha) {
	s = Sha{}

	if err := s.Set(v); err != nil {
		panic(errors.Wrap(err))
	}

	return
}

func MakeSha(v string) (s Sha, err error) {
	s = Sha{}
	err = s.Set(v)

	return
}

func MakeShaFromPath(p string) (s Sha, err error) {
	schwanz := filepath.Base(p)
	kopf := filepath.Base(filepath.Dir(p))

	switch {
	case schwanz == string(filepath.Separator) || kopf == string(filepath.Separator):
		fallthrough

	case schwanz == "." || kopf == ".":
		err = errors.Errorf(
			"path cannot be turned into a kopf/schwanz pair: '%s/%s'",
			kopf,
			schwanz,
		)

		return
	}

	if s, err = MakeSha(fmt.Sprintf("%s%s", kopf, schwanz)); err != nil {
		err = errors.Wrapf(err, "kopf: %q, schwanz: %q", kopf, schwanz)
		return
	}

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
	s.value = fmt.Sprintf("%x", h.Sum(nil))

	return
}

func (s Sha) GetShaString() string {
	return s.String()
}

func (s Sha) String() string {
	if s.value == "" {
		return ShaNull
	} else {
		return s.value
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
		err = errors.Wrapf(err, "%q", v1)
		return
	}

	s.value = v1

	return
}

func (s Sha) GetSha() schnittstellen.Sha {
	return s
}

func (s Sha) IsNull() bool {
	if s.value == "" {
		return true
	}

	if s.value == ShaNull {
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

func (a Sha) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Sha) EqualsSha(b schnittstellen.Sha) bool {
	return a.GetShaString() == b.GetShaString()
}

func (a Sha) Equals(b Sha) bool {
	return a.GetShaString() == b.GetShaString()
}

func (s Sha) Path(pc ...string) string {
	pc = append(pc, s.Kopf())
	pc = append(pc, s.Schwanz())

	return path.Join(pc...)
}

func (s Sha) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Sha) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}

func (s Sha) MarshalText() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Sha) UnmarshalText(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
