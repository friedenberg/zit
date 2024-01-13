package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type (
	Sha struct {
		value sha.Sha
	}

	ShaLike interface {
		GetSha() Sha
	}
)

func MakeShaLike(v string) (t ShaLike, err error) {
	return MakeSha(v)
}

func MakeSha(v string) (t Sha, err error) {
	if t.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MustSha(v string) (t Sha) {
	if err := t.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func (_ Sha) GetQueryPrefix() string {
	return "@"
}

func (t *Sha) Reset() {
	t.value.Reset()
}

func (a *Sha) ResetWith(b Sha) {
	a.value = b.value
}

func (a *Sha) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a *Sha) Equals(b *Sha) bool {
	return a.value.Equals(&b.value)
}

func (e Sha) GetGattung() schnittstellen.GattungLike {
	return gattung.Akte
}

func (e Sha) String() string {
	return e.value.String()
}

func (t Sha) Parts() [3]string {
	return [3]string{"", "@", t.value.String()}
}

func (e *Sha) Set(v string) (err error) {
	v = strings.TrimSpace(strings.Trim(v, "@ "))

	if err = e.value.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Sha) GetSha() Sha {
	return t
}

func (a Sha) EqualsSha(b schnittstellen.ShaLike) bool {
	return a.value.EqualsSha(b)
}

func (t Sha) GetShaBytes() []byte {
	return t.value.GetShaBytes()
}

func (t Sha) GetShaString() string {
	return t.value.GetShaString()
}

func (t Sha) IsNull() bool {
	return t.value.IsNull()
}

func (t Sha) Kopf() string {
	return t.value.Kopf()
}

func (t Sha) Schwanz() string {
	return t.value.Schwanz()
}

func (t Sha) GetShaLike() schnittstellen.ShaLike {
	return t
}

func (t Sha) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Sha) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Sha) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Sha) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
