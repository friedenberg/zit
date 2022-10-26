package int_value

import (
	"strconv"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type IntValue int

func Make(i int) IntValue {
	return IntValue(i)
}

func (a IntValue) Less(b IntValue) (ok bool) {
	ok = a.Int() < b.Int()
	return
}

func (iv *IntValue) Reset() {
	*iv = IntValue(0)
}

func (iv *IntValue) SetInt(i int) {
	*iv = IntValue(i)
}

func (iv *IntValue) Set(v string) (err error) {
	var i int

	if i, err = strconv.Atoi(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	iv.SetInt(i)

	return
}

func (iv IntValue) String() string {
	return strconv.Itoa(int(iv))
}

func (iv IntValue) Int() int {
	return int(iv)
}
