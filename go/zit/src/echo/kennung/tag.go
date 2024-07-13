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
	register(Tag{})
}

const TagRegexString = `^%?[-a-z0-9_]+$`

var TagRegex *regexp.Regexp

func init() {
	TagRegex = regexp.MustCompile(TagRegexString)
}

var (
	sTagResetter  tagResetter
	sTag2Resetter tag2Resetter
)

type Tag = tag

var TagResetter = sTagResetter

// type Etikett = etikett2
// var EtikettResetter = sEtikett2Resetter

type tag struct {
	virtual       bool
	dependentLeaf bool
	value         string
}

type IndexedEtikett = IndexedLike

func MustEtikettPtr(v string) (e *Tag) {
	e = &Tag{}
	e.init()

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MustTag(v string) (e Tag) {
	e.init()

	var err error

	if err = e.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func MakeTag(v string) (e Tag, err error) {
	e.init()

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e tag) init() {
}

func (e *tag) Reset() {
	sTagResetter.Reset(e)
}

func (e tag) GetQueryPrefix() string {
	return "-"
}

func (e tag) IsEmpty() bool {
	return e.value == ""
}

func (e tag) GetGenre() interfaces.Genre {
	return gattung.Etikett
}

func (a tag) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a tag) Equals(b tag) bool {
	return a == b
}

func (e tag) String() string {
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

func (e tag) Bytes() []byte {
	return []byte(e.String())
}

func (e tag) Parts() [3]string {
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

func (e tag) IsVirtual() bool {
	return e.virtual
}

func (e tag) IsDependentLeaf() bool {
	return e.dependentLeaf
}

func (e *tag) TodoSetFromKennung2(v *Id) (err error) {
	return e.Set(v.String())
}

func (e *tag) Set(v string) (err error) {
	v1 := v
	v = strings.ToLower(strings.TrimSpace(v))

	if err = ErrOnConfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !TagRegex.MatchString(v) {
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

func (t tag) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *tag) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t tag) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *tag) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
