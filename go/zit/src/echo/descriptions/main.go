package descriptions

import (
	"bytes"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/query_spec"
)

// TODO-P1 move to catgut.String
type Description struct {
	wasSet bool
	value  string
}

func Make(v string) Description {
	return Description{
		wasSet: true,
		value:  v,
	}
}

func (b Description) String() string {
	return b.value
}

func (b *Description) TodoSetManyCatgutStrings(
	vs ...*catgut.String,
) (err error) {
	var s catgut.String

	if _, err = s.Append(vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return b.Set(s.String())
}

func (b *Description) TodoSetSlice(v catgut.Slice) (err error) {
	return b.Set(v.String())
}

func (b *Description) ReadFromTokenScanner(ts *query_spec.TokenScanner) (err error) {
	var sb strings.Builder

	for ts.Scan() {
		token, tokenType := ts.GetTokenAndType()

		if tokenType == query_spec.TokenTypeOperator &&
			bytes.Equal(token.Bytes(), []byte{'\n'}) {
			ts.Unscan()
			break
		}

		sb.Write(ts.GetToken().Bytes())
	}

	if err = b.Set(sb.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (b *Description) Set(v string) (err error) {
	b.wasSet = true

	v1 := strings.TrimSpace(v)

	if v0 := b.String(); v0 != "" && v0 != v1 {
		b.value = v0 + "\n" + v1
	} else {
		b.value = v1
	}

	return
}

func (a Description) WasSet() bool {
	return a.wasSet
}

func (a *Description) Reset() {
	a.wasSet = false
	a.value = ""
}

func (a Description) IsEmpty() bool {
	return a.value == ""
}

func (a Description) Equals(b Description) (ok bool) {
	// if !a.wasSet {
	// 	return false
	// }

	return a.value == b.value
}

func (a Description) Less(b Description) (ok bool) {
	return a.value < b.value
}

func (a Description) MarshalBinary() (text []byte, err error) {
	text = []byte(a.value)
	return
}

func (a *Description) UnmarshalBinary(text []byte) (err error) {
	a.wasSet = true
	a.value = string(text)
	return
}

func (a Description) MarshalText() (text []byte, err error) {
	text = []byte(a.value)
	return
}

func (a *Description) UnmarshalText(text []byte) (err error) {
	a.wasSet = true
	a.value = string(text)
	return
}
