package kennung

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

const KastenRegexString = `^(//)?[-a-z0-9_]+$`

var KastenRegex *regexp.Regexp

func init() {
	KastenRegex = regexp.MustCompile(KastenRegexString)
	register(Kasten{})
}

func MustKasten(v string) (e Kasten) {
	if err := e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeKasten(v string) (e Kasten, err error) {
	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type Kasten struct {
	value string
}

func (e *Kasten) Reset() {
	e.value = ""
}

func (e *Kasten) ResetWith(e1 Kasten) {
	e.value = e1.value
}

func (a Kasten) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Kasten) Equals(b Kasten) bool {
	return a.value == b.value
}

func (o Kasten) GetGattung() schnittstellen.GattungLike {
	return gattung.Kasten
}

func (k Kasten) String() string {
	return k.value
}

func (k Kasten) Parts() [3]string {
	return [3]string{"/", "/", k.value}
}

func (k Kasten) GetQueryPrefix() string {
	return "//"
}

func (e *Kasten) Set(v string) (err error) {
	v = strings.TrimPrefix(v, "//")
	v = strings.ToLower(strings.TrimSpace(v))

	if err = ErrOnKonfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !KastenRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid Kasten: '%s'", v)
		return
	}

	e.value = v

	return
}

func (t Kasten) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Kasten) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Kasten) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Kasten) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
