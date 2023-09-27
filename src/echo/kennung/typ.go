package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

func init() {
	register(Typ{})
}

type (
	Typ struct {
		value string
	}

	TypLike interface {
		GetTyp() Typ
		GetTypPtr() *Typ
	}

	InlineTypChecker interface {
		IsInlineTyp(Typ) bool
	}
)

func MakeTyp(v string) (t Typ, err error) {
	if t.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MustTyp(v string) (t Typ) {
	if err := t.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func (t Typ) IsEmpty() bool {
	return t.value == ""
}

func (t *Typ) Reset() {
	t.value = ""
}

func (a *Typ) ResetWith(b Typ) {
	a.value = b.value
}

func (a Typ) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Typ) Equals(b Typ) bool {
	return a.value == b.value
}

func (t Typ) GetTyp() Typ {
	return t
}

func (t *Typ) GetTypPtr() *Typ {
	return t
}

func (o Typ) GetGattung() schnittstellen.GattungLike {
	return gattung.Typ
}

func (e Typ) String() string {
	return e.value
}

func (t Typ) Parts() [3]string {
	return [3]string{"", "!", t.value}
}

func (e *Typ) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(strings.Trim(v, ".! ")))

	if err = ErrOnKonfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !EtikettRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid Typ: '%s'", v)
		return
	}

	e.value = v

	return
}

func (t Typ) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Typ) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Typ) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Typ) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
