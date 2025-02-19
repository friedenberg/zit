package bech32

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

type Value struct {
	HRP  string // human-readable part
	Data []byte
}

func (value Value) MarshalText() (text []byte, err error) {
	if text, err = Encode(value.HRP, value.Data); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (value *Value) UnmarshalText(text []byte) (err error) {
	if value.HRP, value.Data, err = Decode(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
