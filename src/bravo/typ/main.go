package typ

import "strings"

var TextTypes map[string]bool

func init() {
	TextTypes = map[string]bool{
		"md": true,
	}
}

type Typ struct {
	Value string
}

func (v Typ) String() string {
	return v.Value
}

func (v *Typ) Set(v1 string) (err error) {
	v.Value = strings.TrimSpace(strings.TrimPrefix(v1, "."))

	return
}

func (v Typ) IsTextType() (is bool) {
	is, _ = TextTypes[v.String()]

	return
}
