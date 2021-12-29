package akte_ext

import "strings"

var TextTypes map[string]bool

func init() {
	TextTypes = map[string]bool{
		"md": true,
	}
}

type AkteExt struct {
	Value string
}

func (v AkteExt) String() string {
	return v.Value
}

func (v *AkteExt) Set(v1 string) (err error) {
	v.Value = strings.TrimSpace(strings.TrimPrefix(v1, "."))

	return
}

func (v AkteExt) IsTextType() (is bool) {
	is, _ = TextTypes[v.String()]

	return
}
