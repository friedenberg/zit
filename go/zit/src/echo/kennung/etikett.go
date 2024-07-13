package kennung

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

func init() {
	register(Etikett{})
}

const EtikettRegexString = `^%?[-a-z0-9_]+$`

var EtikettRegex *regexp.Regexp

func init() {
	EtikettRegex = regexp.MustCompile(EtikettRegexString)
}

var (
	sEtikettResetter  etikettResetter
	sEtikett2Resetter etikett2Resetter
)

type Etikett = etikett

var EtikettResetter = sEtikettResetter

// type Etikett = etikett2
// var EtikettResetter = sEtikett2Resetter

type etikett struct {
	virtual       bool
	dependentLeaf bool
	value         string
}

type IndexedEtikett = IndexedLike

func MustEtikettPtr(v string) (e *Etikett) {
	e = &Etikett{}
	e.init()

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MustEtikett(v string) (e Etikett) {
	e.init()

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeEtikett(v string) (e Etikett, err error) {
	e.init()

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e etikett) init() {
}

func (e *etikett) Reset() {
	sEtikettResetter.Reset(e)
}

func (e etikett) GetQueryPrefix() string {
	return "-"
}

func (e etikett) IsEmpty() bool {
	return e.value == ""
}

func (e etikett) GetGattung() interfaces.GattungLike {
	return gattung.Etikett
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
