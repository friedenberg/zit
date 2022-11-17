package typ

import (
	"strings"

	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/konfig"
)

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

func (t Typ) IsInline(k konfig.Konfig) (isInline bool) {
	ts := t.String()
	isInline = k.Compiled.TypenInline.Contains(ts)

	if typKonfig, ok := k.Typen[ts]; ok {
		isInline = typKonfig.InlineAkte
	}

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
