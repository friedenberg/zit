package int_value

import (
	"strconv"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type IntValue int

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
