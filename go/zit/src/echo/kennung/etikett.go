package kennung

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/gattung"
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

type IndexedEtikett = IndexedLike

func MustEtikettPtr(v string) (e *Etikett) {
	e = &Etikett{}

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
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

func (e Etikett) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (e *Etikett) Reset() {
	EtikettResetter.Reset(e)
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
	v := e.String()

	if strings.HasPrefix(v, "-") {
		v = v[1:]
	}

	return [3]string{"", "-", v}
}

func (e *Etikett) TodoSetFromKennung2(v *Kennung2) (err error) {
	return e.Set(v.String())
}

func (e *Etikett) Set(v string) (err error) {
	v1 := v
	v = strings.ToLower(strings.TrimSpace(v))

	if err = ErrOnKonfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !EtikettRegex.MatchString(v) {
		err = errors.Errorf("not a valid etikett: %q", v1)
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
