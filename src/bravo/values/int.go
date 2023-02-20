package values

import (
	"strconv"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Int int

func MakeInt(i int) Int {
	return Int(i)
}

func (a Int) Equals(b Int) (ok bool) {
	ok = a.Int() == b.Int()

	return
}

func (a Int) Less(b Int) (ok bool) {
	ok = a.Int() < b.Int()
	return
}

func (iv *Int) Reset() {
	*iv = Int(0)
}

func (iv *Int) ResetWith(b Int) {
	iv.SetInt(b.Int())
}

func (iv *Int) SetInt(i int) {
	*iv = Int(i)
}

func (iv *Int) Set(v string) (err error) {
	var i int

	if i, err = strconv.Atoi(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	iv.SetInt(i)

	return
}

func (iv Int) String() string {
	return strconv.Itoa(int(iv))
}

func (iv Int) Int() int {
	return int(iv)
}
