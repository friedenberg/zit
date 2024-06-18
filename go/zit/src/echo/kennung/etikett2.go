package kennung

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
)

const Etikett2RegexString = `^%?[-a-z0-9_]+$`

var Etikett2Regex *regexp.Regexp

func init() {
	Etikett2Regex = regexp.MustCompile(Etikett2RegexString)
}

type etikett2 struct {
	virtual       bool
	dependentLeaf bool
	value         *catgut.String
}

type IndexedEtikett2 = IndexedLike

func (e *etikett2) init() {
	e.value = catgut.GetPool().Get()
}

func (e etikett2) GetQueryPrefix() string {
	return "-"
}

func (e etikett2) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (e *etikett2) Reset() {
	sEtikett2Resetter.Reset(e)
}

func (a etikett2) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a etikett2) Equals(b etikett2) bool {
	return a == b
}

func (e etikett2) String() string {
	var sb strings.Builder

	if e.virtual {
		sb.WriteRune('%')
	}

	if e.dependentLeaf {
		sb.WriteRune('-')
	}

	sb.Write(e.value.Bytes())

	return sb.String()
}

func (e etikett2) Bytes() []byte {
	return []byte(e.String())
}

func (e etikett2) Parts() [3]string {
	switch {
	case e.virtual && e.dependentLeaf:
		return [3]string{"%", "-", e.value.String()}

	case e.virtual:
		return [3]string{"", "%", e.value.String()}

	case e.dependentLeaf:
		return [3]string{"", "-", e.value.String()}

	default:
		return [3]string{"", "", e.value.String()}
	}
}

func (e etikett2) IsVirtual() bool {
	return e.virtual
}

func (e etikett2) IsDependentLeaf() bool {
	return e.dependentLeaf
}

func (e *etikett2) TodoSetFromKennung2(v *Kennung2) (err error) {
	return e.Set(v.String())
}

func (e *etikett2) Set(v string) (err error) {
	v1 := v
	v = strings.ToLower(strings.TrimSpace(v))

	if err = ErrOnKonfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !Etikett2Regex.MatchString(v) {
		err = errors.Errorf("not a valid etikett: %q", v1)
		return
	}

	e.virtual = strings.HasPrefix(v, "%")
	v = strings.TrimPrefix(v, "%")

	e.dependentLeaf = strings.HasPrefix(v, "-")
	v = strings.TrimPrefix(v, "-")

	if e.value == nil {
		e.value = catgut.GetPool().Get()
	}

	if err = e.value.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t etikett2) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *etikett2) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t etikett2) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *etikett2) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
