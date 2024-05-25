package kennung

import (
	"regexp"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/delta/gattung"
)

const Etikett2RegexString = `^%?[-a-z0-9_]+$`

var Etikett2Regex *regexp.Regexp

func init() {
	Etikett2Regex = regexp.MustCompile(Etikett2RegexString)
}

type Etikett2 struct {
	virtual       bool
	dependentLeaf bool
	value         *catgut.String
}

type IndexedEtikett2 = IndexedLike

func MustEtikett2Ptr(v string) (e *Etikett2) {
	e = &Etikett2{
		value: catgut.GetPool().Get(),
	}

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MustEtikett2(v string) (e Etikett2) {
	e.value = catgut.GetPool().Get()
	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeEtikett2(v string) (e Etikett2, err error) {
	e.value = catgut.GetPool().Get()

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e Etikett2) GetQueryPrefix() string {
	return "-"
}

func (e Etikett2) GetGattung() schnittstellen.GattungLike {
	return gattung.Etikett
}

func (e *Etikett2) Reset() {
	Etikett2Resetter.Reset(e)
}

func (a Etikett2) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Etikett2) Equals(b Etikett2) bool {
	return a == b
}

func (e Etikett2) String() string {
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

func (e Etikett2) Bytes() []byte {
	return []byte(e.String())
}

func (e Etikett2) Parts() [3]string {
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

func (e Etikett2) IsVirtual() bool {
	return e.virtual
}

func (e Etikett2) IsDependentLeaf() bool {
	return e.dependentLeaf
}

func (e *Etikett2) TodoSetFromKennung2(v *Kennung2) (err error) {
	return e.Set(v.String())
}

func (e *Etikett2) Set(v string) (err error) {
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

func (t Etikett2) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Etikett2) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Etikett2) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Etikett2) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
