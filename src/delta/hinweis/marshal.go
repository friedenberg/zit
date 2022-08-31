package hinweis

import (
	"bytes"
	"encoding/gob"
)

func (s Hinweis) GobEncode() (b1 []byte, err error) {
  b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
  err = e.Encode(s.inner)
  b1 = b.Bytes()
	return
}

func (s *Hinweis) GobDecode(b []byte) error {
	b1 := bytes.NewBuffer(b)
	d := gob.NewDecoder(b1)
	return d.Decode(&s.inner)
}

func (s Hinweis) MarshalText() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Hinweis) UnmarshalText(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
