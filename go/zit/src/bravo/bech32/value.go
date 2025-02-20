package bech32

import "code.linenisgreat.com/zit/go/zit/src/alfa/errors"

type Value struct {
	HRP  string // human-readable part
	Data []byte
}

func (value Value) String() string {
	var text []byte
	var err error

	if text, err = Encode(value.HRP, value.Data); err != nil {
		panic(err)
	}

	return string(text)
}

func (value Value) MarshalText() (text []byte, err error) {
	if len(value.Data) == 0 {
		return
	}

	if text, err = Encode(value.HRP, value.Data); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (value *Value) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		return
	}

	if value.HRP, value.Data, err = Decode(string(text)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
