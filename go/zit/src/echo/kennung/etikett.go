package kennung

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/delta/gattung"
)

func init() {
	register(Etikett{})
}

const EtikettRegexString = `^%?[-a-z0-9_]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

type Etikett = etikett

var (
	EtikettResetter  etikettResetter
	Etikett2Resetter etikett2Resetter
)

type etikett struct {
	virtual       bool
	dependentLeaf bool
	value         string
}

type IndexedEtikett = IndexedLike

func MustEtikettPtr(v string) (e *etikett) {
	e = &etikett{}

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MustEtikett(v string) (e etikett) {
	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeEtikett(v string) (e etikett, err error) {
	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e etikett) GetQueryPrefix() string {
	return "-"
}

func (e etikett) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (e *etikett) Reset() {
	EtikettResetter.Reset(e)
}

func (a etikett) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a etikett) Equals(b etikett) bool {
	return a == b
}

func (e etikett) String() string {
	var sb strings.Builder

	if e.virtual {
		sb.WriteRune('%')
	}

	if e.dependentLeaf {
		sb.WriteRune('-')
	}

	sb.WriteString(e.value)

	return sb.String()
}

func (e etikett) Bytes() []byte {
	return []byte(e.String())
}

func (e etikett) Parts() [3]string {
	switch {
	case e.virtual && e.dependentLeaf:
		return [3]string{"%", "-", e.value}

	case e.virtual:
		return [3]string{"", "%", e.value}

	case e.dependentLeaf:
		return [3]string{"", "-", e.value}

	default:
		return [3]string{"", "", e.value}
	}
}

func (e etikett) IsVirtual() bool {
	return e.virtual
}

func (e etikett) IsDependentLeaf() bool {
	return e.dependentLeaf
}

func (e *etikett) TodoSetFromKennung2(v *Kennung2) (err error) {
	return e.Set(v.String())
}

func (e *etikett) Set(v string) (err error) {
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

	e.virtual = strings.HasPrefix(v, "%")
	v = strings.TrimPrefix(v, "%")

	e.dependentLeaf = strings.HasPrefix(v, "-")
	v = strings.TrimPrefix(v, "-")

	e.value = v

	return
}

func (t etikett) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *etikett) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t etikett) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *etikett) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
