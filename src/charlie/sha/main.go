package sha

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"path"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/catgut"
)

const (
	ByteSize      = 32
	ShaNullString = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	Null          = ShaNullString
)

var shaNull Sha

func init() {
	errors.PanicIfError(shaNull.Set(ShaNullString))
}

type PathComponents interface {
	PathComponents() []string
}

type ShaLike = schnittstellen.ShaGetter

type Sha struct {
	data catgut.String
}

func (s Sha) GetShaBytes() []byte {
	if s.IsNull() {
		return shaNull.data.Bytes()
	} else {
		return s.data.Bytes()
	}
}

func (s Sha) GetShaString() string {
	if s.IsNull() {
		return fmt.Sprintf("%x", shaNull.data.Bytes())
	} else {
		return fmt.Sprintf("%x", s.data.Bytes())
	}
}

func (s Sha) String() string {
	return s.GetShaString()
}

func (s Sha) Sha() Sha {
	return s
}

func (dst *Sha) SetShaLike(src schnittstellen.ShaLike) (err error) {
	err = dst.data.SetBytes(src.GetShaBytes())
	return
}

func (s *Sha) SetParts(a, b string) (err error) {
	if err = s.Set(a + b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Sha) SetHexBytes(b []byte) (err error) {
	s.Reset()

	b = bytes.TrimSpace(b)

	if len(b) == 0 {
		return
	}

	b1 := s.data.AvailableBuffer()

	var n int

	if b1, n, err = hexDecode(b1, b); err != nil {
		err = errors.Wrapf(err, "N: %d, Data: %q", n, b)
		return
	}

	s.data.Write(b1)

	return
}

func (s *Sha) Set(v string) (err error) {
	s.Reset()

	v1 := strings.TrimSpace(v)

	var b []byte

	if b, err = hex.DecodeString(v1); err != nil {
		err = errors.Wrapf(err, "%q", v1)
		return
	}

	if err = makeErrLength(ByteSize, len(b)); err != nil {
		return
	}

	s.data.Write(b)

	return
}

func (s Sha) GetShaLike() schnittstellen.ShaLike {
	return s
}

func (s Sha) IsNull() bool {
	if s.data.Len() == 0 {
		return true
	}

	if bytes.Equal(s.data.Bytes(), shaNull.data.Bytes()) {
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

func (a Sha) EqualsSha(b schnittstellen.ShaLike) bool {
	return a.GetShaString() == b.GetShaString()
}

func (a Sha) Equals(b Sha) bool {
	return a.GetShaString() == b.GetShaString()
}

func (s *Sha) Reset() {
	s.data.Reset()
	s.data.Grow(ByteSize)
}

func (a *Sha) ResetWith(b Sha) {
	a.data.Reset()
	a.data.Grow(ByteSize)
	errors.PanicIfError(b.data.CopyTo(&a.data))
}

func (a *Sha) ResetWithShaLike(b schnittstellen.ShaLike) {
	a.data.Reset()
	a.data.Grow(ByteSize)
	a.data.Write(b.GetShaBytes())
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
