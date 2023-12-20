package sha

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"path"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
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
	data *[ByteSize]byte
}

func (s *Sha) GetShaBytes() []byte {
	if s.IsNull() {
		return shaNull.data[:]
	} else {
		return s.data[:]
	}
}

func (s *Sha) GetShaString() string {
	if s == nil || s.IsNull() {
		return fmt.Sprintf("%x", shaNull.data[:])
	} else {
		return fmt.Sprintf("%x", s.data[:])
	}
}

func (s *Sha) String() string {
	return s.GetShaString()
}

func (s *Sha) Sha() *Sha {
	return s
}

func (dst *Sha) SetFromHash(h hash.Hash) (err error) {
  dst.allocDataIfNecessary()
	b := h.Sum(dst.data[:0])
	err = makeErrLength(ByteSize, len(b))
	return
}

func (dst *Sha) SetShaLike(src ShaLike) (err error) {
  dst.allocDataIfNecessary()

	err = makeErrLength(
		ByteSize,
		copy(dst.data[:], src.GetShaLike().GetShaBytes()),
	)

	return
}

func (s *Sha) SetParts(a, b string) (err error) {
	if err = s.Set(a + b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (src *Sha) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = w.Write(src.GetShaBytes())
	n = int64(n1)
	return
}

func (s *Sha) ReadFrom(r io.Reader) (n int64, err error) {
  s.allocDataIfNecessary()

	var n1 int

	n1, err = r.Read(s.data[:])

	if n1 == 0 && err == io.EOF {
		return
	} else if n1 != ByteSize && n1 != 0 {
		err = errors.Errorf("expected to read %d bytes but only read %d", ByteSize, n1)
		return
	} else if errors.IsNotNilAndNotEOF(err) {
		err = errors.Wrap(err)
		return
	}

	n = int64(n1)

	return
}

func (s *Sha) SetHexBytes(b []byte) (err error) {
  s.allocDataIfNecessary()

	b = bytes.TrimSpace(b)

	if len(b) == 0 {
		return
	}

	var n int

	if _, n, err = hexDecode(s.data[:0], b); err != nil {
		err = errors.Wrapf(err, "N: %d, Data: %q", n, b)
		return
	}

	return
}

func (s *Sha) Set(v string) (err error) {
  s.allocDataIfNecessary()

	v1 := strings.TrimSpace(v)

	var b []byte

	if b, err = hex.DecodeString(v1); err != nil {
		err = errors.Wrapf(err, "%q", v1)
		return
	}

	n := copy(s.data[:], b)

	if err = makeErrLength(ByteSize, n); err != nil {
		return
	}

	return
}

func (s *Sha) GetShaLike() schnittstellen.ShaLike {
	return s
}

func (s *Sha) IsNull() bool {
	if s == nil || s.data == nil {
		return true
	}

	if bytes.Equal(s.data[:], shaNull.data[:]) {
		return true
	}

	return false
}

func (s *Sha) Kopf() string {
	return s.String()[0:2]
}

func (s *Sha) Schwanz() string {
	return s.String()[2:]
}

func (a *Sha) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a *Sha) EqualsSha(b schnittstellen.ShaLike) bool {
	return a.GetShaString() == b.GetShaString()
}

func (a *Sha) Equals(b *Sha) bool {
	return a.GetShaString() == b.GetShaString()
}


func (s *Sha) allocDataIfNecessary() {
  if s.data != nil {
    return
  }

	s.data = &[ByteSize]byte{}
}

func (s *Sha) Reset() {
  s.allocDataIfNecessary()
	s.ResetWith(&shaNull)
}

func (a *Sha) ResetWith(b *Sha) {
  a.allocDataIfNecessary()
	copy(a.data[:], b.data[:])
}

func (a *Sha) ResetWithShaLike(b schnittstellen.ShaLike) {
	copy(a.data[:], b.GetShaBytes())
}

func (s *Sha) Path(pc ...string) string {
	pc = append(pc, s.Kopf())
	pc = append(pc, s.Schwanz())

	return path.Join(pc...)
}

func (s *Sha) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Sha) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}

func (s *Sha) MarshalText() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Sha) UnmarshalText(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
