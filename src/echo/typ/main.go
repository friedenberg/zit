package typ

import (
	"strings"

	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/delta/konfig"
)

type Kennung struct {
	etikett.Etikett
}

func Make(v string) Kennung {
	return Kennung{
		Etikett: etikett.Etikett{
			Value: v,
		},
	}
}

func (v *Kennung) Set(v1 string) (err error) {
	return v.Etikett.Set(strings.TrimSpace(strings.Trim(v1, ".! ")))
}

func (v Kennung) Expanded() Set {
	return ExpanderRight.Expand(v.String())
}

func (t Kennung) IsInlineAkte(k konfig.Konfig) (isInline bool) {
	ts := t.String()
	tc := k.GetTyp(ts)

	if tc == nil {
		return
	}

	isInline = tc.InlineAkte

	return
}

func (t Kennung) MarshalText() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *Kennung) UnmarshalText(text []byte) (err error) {
	if err = t.Set(string(text)); err != nil {
		return
	}

	return
}

func (t Kennung) MarshalBinary() (text []byte, err error) {
	text = []byte(t.String())

	return
}

func (t *Kennung) UnmarshalBinary(text []byte) (err error) {
	t.Etikett.Value = string(text)

	return
}
