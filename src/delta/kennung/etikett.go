package kennung

import (
	"regexp"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/values"
)

func init() {
	register(Etikett{})
}

const EtikettRegexString = `^[-a-z0-9_]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

type Etikett struct {
	value string
}

func MustEtikett(v string) (e Etikett) {
	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeEtikett(v string) (e Etikett, err error) {
	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e Etikett) GetQueryPrefix() string {
	return "-"
}

func (e Etikett) GetGattung() schnittstellen.Gattung {
	return gattung.Etikett
}

func (e *Etikett) ResetWith(e1 Etikett) {
	*e = e1
}

func (e *Etikett) Reset() {
	e.value = ""
}

func (a Etikett) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Etikett) Equals(b Etikett) bool {
	return a == b
}

func (e Etikett) String() string {
	return e.value
}

func (e Etikett) Parts() [3]string {
	v := e.value

	if strings.HasPrefix(v, "-") {
		v = v[1:]
	}

	return [3]string{"", "-", v}
}

func (e Etikett) ContainsMatchableExactly(m Matchable) bool {
	es := m.GetEtikettenExpanded()

	if es.Contains(e) {
		return true
	}

	e1, ok := m.GetIdLike().(Etikett)

	if ok && e.Equals(e1) {
		return true
	}

	return false
}

func (e Etikett) ContainsMatchable(m Matchable) bool {
	es := m.GetEtikettenExpanded()

	if es.Contains(e) {
		return true
	}

	e1, ok := m.GetIdLike().(Etikett)

	if ok && ContainsWithoutUnderscoreSuffix(e1, e) {
		return true
	}

	return false
}

func (e *Etikett) Set(v string) (err error) {
	if !EtikettRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid etikett: '%s'", v)
		return
	}

	e.value = v

	return
}

func (t Etikett) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Etikett) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Etikett) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Etikett) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Etikett) KennungSansGattungClone() KennungSansGattung {
	return t
}

func (t Etikett) KennungSansGattungPtrClone() KennungSansGattungPtr {
	return &t
}

func (t Etikett) KennungClone() Kennung {
	return t
}

func (t Etikett) KennungPtrClone() KennungPtr {
	return &t
}
