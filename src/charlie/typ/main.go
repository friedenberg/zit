package typ

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
)

var TextTypes map[string]bool

func init() {
	TextTypes = map[string]bool{
		"md": true,
	}
}

type Typ struct {
	typ string
}

func (v Typ) String() string {
	return v.typ
}

func Make(v string) Typ {
	return Typ{
		typ: v,
	}
}

func (v *Typ) Set(v1 string) (err error) {
	v.typ = strings.TrimSpace(strings.TrimPrefix(v1, "."))

	return
}

func (v Typ) IsTextType() (is bool) {
	is, _ = TextTypes[v.String()]

	return
}

func (t Typ) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(t.typ)

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
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
	err = e.Encode(s.typ)
	b1 = b.Bytes()
	return
}

func (s *Typ) GobDecode(b []byte) error {
	b1 := bytes.NewBuffer(b)
	d := gob.NewDecoder(b1)
	return d.Decode(&s.typ)
}
