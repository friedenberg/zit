package kennung

import (
	"regexp"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

func init() {
	register(Etikett{})
}

type EtikettLike interface {
	GetEtikett() Etikett
	GetEtikettPtr() *Etikett
	schnittstellen.StringerGattungGetter
}

const EtikettRegexString = `^[-a-z0-9_]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

type Etikett struct {
	value string
}

type IndexedEtikett = IndexedLike[Etikett, *Etikett]

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

func (e *Etikett) GetEtikett() Etikett {
	return *e
}

func (e *Etikett) GetEtikettPtr() *Etikett {
	return e
}

func (e Etikett) GetGattung() schnittstellen.GattungLike {
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

func (e Etikett) Bytes() []byte {
	return []byte(e.value)
}

func (e Etikett) Parts() [3]string {
	v := e.value

	if strings.HasPrefix(v, "-") {
		v = v[1:]
	}

	return [3]string{"", "-", v}
}

func (e *Etikett) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(v))

	if err = ErrOnKonfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

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
