package ids

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

func init() {
	register(Type{})
}

type (
	Type struct {
		value string
	}

	InlineTypeChecker interface {
		IsInlineType(Type) bool
	}
)

func MakeType(v string) (t Type, err error) {
	if t.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MustType(v string) (t Type) {
	if err := t.Set(v); err != nil {
		errors.PanicIfError(err)
	}

	return
}

func (t Type) IsEmpty() bool {
	return t.value == ""
}

func (t *Type) Reset() {
	t.value = ""
}

func (a *Type) ResetWith(b Type) {
	a.value = b.value
}

func (a Type) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Type) Equals(b Type) bool {
	return a.value == b.value
}

func (t Type) GetType() Type {
	return t
}

func (t *Type) GetTypPtr() *Type {
	return t
}

func (o Type) GetGenre() interfaces.Genre {
	return genres.Type
}

func (t Type) IsToml() bool {
	return strings.HasPrefix(t.value, "toml")
}

func (e Type) String() string {
	return e.value
}

func (t Type) Parts() [3]string {
	return [3]string{"", "!", t.value}
}

func (e *Type) TodoSetFromObjectId(v *ObjectId) (err error) {
	return e.Set(v.String())
}

func (e *Type) Set(v string) (err error) {
	v = strings.ToLower(strings.TrimSpace(strings.Trim(v, ".! ")))

	if err = ErrOnConfig(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !TagRegex.Match([]byte(v)) {
		err = errors.Errorf("not a valid Typ: '%s'", v)
		return
	}

	e.value = v

	return
}

func (t Type) MarshalText() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Type) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (t Type) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())
	return
}

func (t *Type) UnmarshalBinary(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
