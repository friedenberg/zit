package sha

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"path"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
)

const (
	ByteSize      = 32
	ShaNullString = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	Null          = ShaNullString
)

type Bytes [ByteSize]byte

var shaNull Sha

func init() {
	errors.PanicIfError(shaNull.Set(ShaNullString))
}

type PathComponents interface {
	PathComponents() []string
}

type ShaLike = interfaces.ShaGetter

type Sha struct {
	data *Bytes
}

func (s *Sha) Size() int {
	return ByteSize
}

func (s *Sha) GetBytes() Bytes {
	if s.IsNull() {
		return *shaNull.data
	} else {
		return *s.data
	}
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

func (src *Sha) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = ohio.WriteAllOrDieTrying(w, src.GetShaBytes())
	n = int64(n1)
	return
}

func (s *Sha) GetShaLike() interfaces.Sha {
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

func (s *Sha) GetHead() string {
	return s.String()[0:2]
}

func (s *Sha) GetTail() string {
	return s.String()[2:]
}

func (a *Sha) AssertEqualsShaLike(b interfaces.Sha) error {
	if !a.EqualsSha(b) {
		return MakeErrNotEqual(a, b)
	}

	return nil
}

func (a *Sha) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a *Sha) EqualsSha(b interfaces.Sha) bool {
	return a.GetShaString() == b.GetShaString()
}

func (a *Sha) Equals(b *Sha) bool {
	return a.GetShaString() == b.GetShaString()
}

//  __        __    _ _   _
//  \ \      / / __(_) |_(_)_ __   __ _
//   \ \ /\ / / '__| | __| | '_ \ / _` |
//    \ V  V /| |  | | |_| | | | | (_| |
//     \_/\_/ |_|  |_|\__|_|_| |_|\__, |
//                                |___/

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

	// if dst.String() == "f32cf7f2b1b8b7688c78dec4eb3c3675fcdc3aeaaf4c1305f3bdaf7fc0252e02" {
	// 	panic("found")
	// }

	return
}

func (s *Sha) SetParts(a, b string) (err error) {
	if err = s.Set(a + b); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Sha) ReadAtFrom(r io.ReaderAt, start int64) (n int64, err error) {
	s.allocDataIfNecessary()

	var n1 int
	n1, err = r.ReadAt(s.data[:], start)
	n += int64(n1)

	if n == 0 && err == io.EOF {
		return
	} else if n != ByteSize && n != 0 {
		err = errors.Errorf("expected to read %d bytes but only read %d", ByteSize, n)
		return
	} else if errors.IsNotNilAndNotEOF(err) {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Sha) ReadFrom(r io.Reader) (n int64, err error) {
	s.allocDataIfNecessary()

	var n1 int
	n1, err = ohio.ReadAllOrDieTrying(r, s.data[:])
	n += int64(n1)

	if n == 0 && err == io.EOF {
		return
	} else if n != ByteSize && n != 0 {
		err = errors.Errorf("expected to read %d bytes but only read %d", ByteSize, n)
		return
	} else if errors.IsNotNilAndNotEOF(err) {
		err = errors.Wrap(err)
		return
	}

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
	v1 = strings.TrimPrefix(v, "@")

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

func (s *Sha) allocDataIfNecessary() {
	if s.data != nil {
		return
	}

	s.data = &Bytes{}
}

func (s *Sha) Reset() {
	s.allocDataIfNecessary()
	s.ResetWith(&shaNull)
}

func (a *Sha) ResetWith(b *Sha) {
	a.allocDataIfNecessary()

	if b.IsNull() {
		copy(a.data[:], shaNull.data[:])
	} else {
		copy(a.data[:], b.data[:])
	}
}

func (a *Sha) ResetWithShaLike(b interfaces.Sha) {
	a.allocDataIfNecessary()
	copy(a.data[:], b.GetShaBytes())
}

func (s *Sha) Path(pc ...string) string {
	pc = append(pc, s.GetHead())
	pc = append(pc, s.GetTail())

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
