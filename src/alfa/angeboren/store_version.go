package angeboren

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/int_value"
)

type storeVersion int_value.IntValue

func (a storeVersion) Less(b schnittstellen.StoreVersion) bool {
	return a.String() < b.String()
}

func (a storeVersion) String() string {
	return int_value.IntValue(a).String()
}

func (a storeVersion) Int() int {
	return int_value.IntValue(a).Int()
}
