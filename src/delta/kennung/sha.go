package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/values"
)

type (
	Sha struct {
		value sha.Sha
	}

	ShaLike interface {
		GetSha() sha.Sha
		Matcher
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

func (a Sha) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Sha) Equals(b Sha) bool {
	return a.value.Equals(b.value)
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

func (t Sha) Each(_ schnittstellen.FuncIter[Matcher]) error {
	return nil
}

func (t Sha) MatcherLen() int {
	return 0
}

func (t Sha) ContainsMatchableExactly(m Matchable) bool {
	return t.ContainsMatchable(m)
}

func (t Sha) ContainsMatchable(m Matchable) bool {
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

func (t Sha) KennungClone() Kennung {
	return t
}

func (t Sha) KennungPtrClone() KennungPtr {
	return &t
}

func (t Sha) KennungSansGattungClone() KennungSansGattung {
	return t
}

func (t Sha) KennungSansGattungPtrClone() KennungSansGattungPtr {
	return &t
}
