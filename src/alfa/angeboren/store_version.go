package angeboren

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/values"
)

type storeVersion values.Int

func (a storeVersion) Less(b schnittstellen.StoreVersion) bool {
	return a.String() < b.String()
}

func (a storeVersion) String() string {
	return values.Int(a).String()
}

func (a storeVersion) Int() int {
	return values.Int(a).Int()
}
