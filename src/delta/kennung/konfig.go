package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/values"
)

func ErrOnKonfig(v string) (err error) {
	if v == "konfig" {
		return errors.Errorf("cannot be %q", "konfig")
	}

	return nil
}

type Konfig struct{}

func (a Konfig) GetGattung() schnittstellen.GattungLike {
	return gattung.Konfig
}

func (a Konfig) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Konfig) Equals(b Konfig) bool {
	return true
}

func (a *Konfig) Reset() {
	return
}

func (a *Konfig) ResetWith(_ Konfig) {
	return
}

func (i Konfig) String() string {
	return "konfig"
}

func (i Konfig) ContainsMatchable(m Matchable) bool {
	if !i.GetGattung().EqualsGattung(gattung.Make(m.GetGattung())) {
		return false
	}

	return true
}

func (k Konfig) Parts() [3]string {
	return [3]string{"", "", "konfig"}
}

func (i Konfig) Set(v string) (err error) {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	if v != "konfig" {
		err = errors.Errorf("not konfig")
		return
	}

	return
}

func (t Konfig) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Konfig) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Konfig) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Konfig) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Konfig) KennungSansGattungClone() KennungSansGattung {
	return t
}

func (t Konfig) KennungSansGattungPtrClone() KennungSansGattungPtr {
	return &t
}

func (t Konfig) KennungClone() Kennung {
	return t
}

func (t Konfig) KennungPtrClone() KennungPtr {
	return &t
}
