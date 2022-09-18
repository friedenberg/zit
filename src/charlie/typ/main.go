package typ

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/friedenberg/zit/src/charlie/etikett"
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

func (s Typ) GobEncode() (b1 []byte, err error) {
	b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
	err = e.Encode(s.String())
	b1 = b.Bytes()
	return
}

func (s *Typ) GobDecode(b []byte) error {
	b1 := bytes.NewBuffer(b)
	d := gob.NewDecoder(b1)
	return d.Decode(&s.Etikett.Value)
}
