package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
)

type (
	Sha struct {
		value sha.Sha
	}
)

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

func (t *Sha) Reset() {
	t.value.Reset()
}

func (a *Sha) ResetWith(b Sha) {
	a.value = b.value
}

func (a Sha) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Sha) Equals(b Sha) bool {
	return a.value.Equals(b.value)
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

func (t Sha) ContainsMatchable(m Matchable) bool {
	if t.value.EqualsSha(m.GetObjekteSha()) {
		return true
	}

	if t.value.EqualsSha(m.GetAkteSha()) {
		return true
	}

	return false
}

func (t Sha) GetSha() sha.Sha {
  return t.value
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

func (t Sha) KennungSansGattungClone() KennungSansGattung {
	return t
}
