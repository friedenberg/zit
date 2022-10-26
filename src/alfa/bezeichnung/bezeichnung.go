package bezeichnung

import "strings"

type Bezeichnung string

func (b Bezeichnung) String() string {
	return string(b)
}

func (b *Bezeichnung) Set(v string) (err error) {
	v1 := strings.TrimSpace(v)

	if v0 := b.String(); v0 != "" {
		*b = Bezeichnung(b.String() + " " + v1)
	} else {
		*b = Bezeichnung(v1)
	}

	return
}

func (a Bezeichnung) Equals(b Bezeichnung) (ok bool) {
	return a == b
}
