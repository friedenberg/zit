package typ

import (
	"strings"

	"github.com/friedenberg/zit/src/delta/etikett"
)

var TextTypes map[string]bool

func init() {
	TextTypes = map[string]bool{
		"md": true,
	}
}

type Typ struct {
	etikett.Etikett
}

func Make(v string) Typ {
	return Typ{
		Etikett: etikett.Etikett{
			Value: v,
		},
	}
}

func (v *Typ) Set(v1 string) (err error) {
	return v.Etikett.Set(strings.TrimSpace(strings.Trim(v1, ".! ")))
}

func (v Typ) IsTextType() (is bool) {
	is, _ = TextTypes[v.String()]

	return
}

func (t Typ) MarshalText() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *Typ) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}

func (t Typ) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *Typ) UnmarshalBinary(text []byte) (err error) {
	t.Etikett.Value = string(text)

	return
}
