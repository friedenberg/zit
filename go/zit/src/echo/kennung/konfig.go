package kennung

import (
	"bytes"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

func init() {
	register(Konfig{})
}

var konfigBytes = []byte("konfig")

func ErrOnKonfigBytes(b []byte) (err error) {
	if bytes.Equal(b, konfigBytes) {
		return errors.Errorf("cannot be %q", "konfig")
	}

	return nil
}

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
